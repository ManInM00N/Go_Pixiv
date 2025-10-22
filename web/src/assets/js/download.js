import axios from "axios";
import {DownloadByFollowPage} from "../../../bindings/main/internal/pixivlib/ctl.js";
import JSZip from "jszip"
import {ElNotification} from "element-plus";
import {sleep} from "./Time.js";
function Download() {
    axios.post("http://127.0.0.1:7234/api/download", {
        type: "Pid",

    }, {
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(response => {

    }).catch(error => {
        console.error
    })
}
export function DownloadFollow(from ,to){
    for (let i =from ; i<=to; i++){
        DownloadByFollowPage(i.toString(), "all")
    }
}


class WorkerManager {
    constructor(config = {}) {
        // 配置
        this.maxConcurrent = config.maxConcurrent || 3;  // 同时最多3个GIF（即3个worker）
        this.workersPerGif = 1;  // 每个GIF固定用1个worker
        this.defaultQuality = config.defaultQuality || 10;

        // 队列管理
        this.queue = [];           // 等待队列（无上限）
        this.processing = [];      // 正在处理的任务
        this.completed = 0;        // 已完成数量
        this.failed = 0;           // 失败数量
        this.Running = false
    }
    Stop(){
        this.Running = false
    }
    async Run(){
        this.Running = true
        try{
            for(;this.Running;){
                await this.processQueue();
                await sleep(1000)
                console.log('[Queue Status]', {
                    queue: this.queue.length,
                    processing: this.processing.length,
                    completed: this.completed,
                    failed: this.failed
                });
            }
        }catch (e){
            console.log(e)
        }
    }

    /**
     * 创建 GIF（自动加入队列）
     * @param {Object} options - GIF选项
     * @returns {Promise<Blob>}
     */
    async createGIF(options = {}) {
        return new Promise((resolve, reject) => {
            const task = {
                options: {
                    workers: this.workersPerGif,  // 固定1个worker
                    quality: options.quality || this.defaultQuality,
                    width: options.width,
                    height: options.height,
                    frames: options.frames || [],
                    pid: options.pid || 'unknown',
                    url: options.url,
                    identify: options.identify
                },
                resolve,
                reject,
                id: `${options.pid || Date.now()}-${Math.random()}`,
                createdAt: Date.now(),
                status:"pending"
            };

            this.queue.push(task);
            console.log(`[GIF Queue] Task ${task.options.pid} added. Queue: ${this.queue.length}`);

        });
    }

    /**
     * 处理队列
     */
    async processQueue() {
        // 如果已达到并发上限，不处理
        if (this.processing.length >= this.maxConcurrent) {
            return;
        }

        // 如果队列为空，不处理
        if (this.queue.length === 0) {
            return;
        }

        // 从队列中取出任务
        const task = this.queue.shift();
        this.processing.push(task);

        console.log(`[GIF Queue] Processing PID ${task.options.pid}. Queue: ${this.queue.length}, Processing: ${this.processing.length}`);

        // 异步处理任务，不阻塞队列
        this.processTaskAsync(task);

    }

    /**
     * 异步处理任务（不阻塞队列）
     */
    async processTaskAsync(task) {
        try {
            const blob = await this.processTask(task);

            // 任务完成
            this.removeFromProcessing(task);
            this.completed++;
            task.resolve(blob);

            console.log(`[GIF Queue] PID ${task.options.pid} completed. Total: ${this.completed}`);
        } catch (error) {
            // 任务失败
            this.removeFromProcessing(task);
            this.failed++;
            task.reject(error);

            console.error(`[GIF Queue] PID ${task.options.pid} failed:`, error);
        }
    }

    /**
     * 实际处理单个 GIF 任务
     */
    async processTask(task) {
        let gif = null

        let zip = null;
        let imgElements = [];
        try {
            let response = await axios.get("http://127.0.0.1:7234/api/getugoira?url=" + task.options.url, {responseType: 'blob'});
            zip = await JSZip.loadAsync(response.data);
            const images = [];
            task.status = "loading frames"
            for (const filename of Object.keys(zip.files)) {
                if (/\.(jpg|jpeg|png|gif)$/i.test(filename)) {
                    const file = zip.files[filename];
                    const imgBlob = await file.async('blob');
                    const imgUrl = URL.createObjectURL(imgBlob);
                    const img = new Image();
                    img.src = imgUrl;
                    images.push(new Promise(resolve => {
                        img.onload = () => {
                            resolve(img);
                            task.options.width = img.width
                            task.options.height = img.height
                        };
                        // 添加错误处理
                        img.onerror = () => {
                            console.error(`Failed to load image: ${filename}`);
                            resolve(null);
                        };
                    }));
                }
            }
            imgElements = await Promise.all(images);
            gif = new window.GIF({
                workers: task.options.workers,
                quality: 1,
                width: task.options.width,
                height: task.options.height,
                worker: '/gif.worker.js',
            });
            imgElements.forEach((img, index) => {
                const delay = task.options.frames[index].delay || 500;  // 默认延迟 500ms
                gif.addFrame(img, {delay: delay});
            });
            task.status = "rendering"

            const blob = await new Promise((resolve, reject) => {
                let isResolved = false;

                const timeoutId = setTimeout(() => {
                    if (!isResolved) {
                        isResolved = true;
                        gif.abort();
                        reject(new Error(`GIF rendering timeout for ${task.options.pid}`));
                    }
                }, 180000); // 180秒超时

                gif.on('finished',  (blob) => {
                    if (!isResolved) {
                        isResolved = true;
                        clearTimeout(timeoutId);
                        resolve(blob)
                    }
                });
                gif.on('abort', (error) => {
                    if (!isResolved) {
                        isResolved = true;
                        clearTimeout(timeoutId);
                        reject(error);
                    }
                });
                gif.render();
            });
            task.status = "uploading";
            await uploadZipAndGenerateGIF(blob, task.options.pid,task.options.identify);

            task.status = "finished";
            console.log(`[${task.options.pid}] All done!`);

            imgElements.forEach(img => {
                if (img && img.src) {
                    URL.revokeObjectURL(img.src);
                }
            });

            return blob;
        } catch (error) {
            task.status = "failed";
            console.error("Error processing GIF:", error);
            ElNotification({
                type: 'error',
                message: 'Error processing GIF: ' + task.options.pid + ' ' + error.message,
                title: "Download GIF failed",
                duration:10000,
            });
            throw error;
        } finally {
            if (gif) {
                try {
                    gif.abort();
                } catch (e) {
                    // ignore
                }
            }
        }
    }

    /**
     * 从处理列表中移除任务
     */
    removeFromProcessing(task) {
        const index = this.processing.findIndex(t => t.id === task.id);
        if (index !== -1) {
            this.processing.splice(index, 1);
            console.log(`[GIF Queue] Removed task ${task.options.pid} from processing. Remaining: ${this.processing.length}`);
        } else {
            console.warn(`[GIF Queue] Task ${task.options.pid} not found in processing list`);
        }
    }

    /**
     * 取消所有等待中的任务
     */
    clearQueue() {
        const canceledCount = this.queue.length;
        this.queue.forEach(task => {
            task.reject(new Error('Task canceled'));
        });
        this.queue = [];
        console.log(`[GIF Queue] Cleared ${canceledCount} pending tasks`);
    }

    /**
     * 取消所有任务
     */
    cancelAll() {
        this.processing.forEach(task => {
            task.status = "canceled";
            task.reject(new Error('Task canceled'));
        });
        this.processing = [];
        this.clearQueue();
        console.log(`[GIF Queue] Canceled all tasks`);
    }

    /**
     * 获取队列状态
     */
    getStatus() {
        return {
            waiting: this.queue.length,
            processing: this.processing.length,
            completed: this.completed,
            failed: this.failed,
            totalWorkers: this.processing.length,
            maxWorkers: this.maxConcurrent,
            tasks: this.processing.map(t => ({
                pid: t.options.pid,
                status: t.status
            }))
        };
    }

    /**
     * 设置最大并发数
     */
    setMaxConcurrent(max) {
        this.maxConcurrent = max;
    }
}

const workerManager = new WorkerManager({
    maxConcurrent: 3,      // 同时最多3个GIF = 3个worker
    defaultQuality: 10     // 默认质量
});
workerManager.Run()
/**
 * 下载并生成 GIF
 * @param {string} pid - 作品ID
 * @param {number} width - 宽度
 * @param {number} height - 高度
 * @param {Array} frames - 帧信息数组（包含delay等）
 * @param {string} url - 下载URL
 * @param {string} identify - 任务标识
 */
export async function DownloadGIF(pid,width,height,frames,url,identify){
    let conf = {
        pid:pid,
        width:width,
        height:height,
        frames:frames,
        url:url,
        identify:identify,
    }
    await workerManager.createGIF(conf)
}
async function uploadZipAndGenerateGIF(zipFile,pid,identify) {
    const formData = new FormData();
    formData.append('file', zipFile);
    try {
        const response = await axios.post('http://127.0.0.1:7234/api/saveugoira?id='+pid+'&identify='+identify, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
    } catch (error) {
        ElNotification({
            type: 'error',
            message: 'Error uploading GIF: ' +pid + ' '+ error,
            title:"Download GIF failed "
        })
        console.error('Error uploading GIF: ' +pid + ' '+ error);
    }
}


/**
 * 批量下载 GIF（示例用法）
 */
export async function batchDownloadGIF(gifList) {
    console.log(`Starting batch download: ${gifList.length} items`);

    const promises = gifList.map(item => {
        return DownloadGIF(
            item.pid,
            item.width,
            item.height,
            item.frames,
            item.url
        ).catch(error => {
            console.error(`Failed to download ${item.pid}:`, error);
            return null;  // 继续处理其他任务
        });
    });

    const results = await Promise.allSettled(promises);

    const succeeded = results.filter(r => r.status === 'fulfilled' && r.value !== null).length;
    const failed = results.length - succeeded;

    console.log(`Batch complete: ${succeeded} succeeded, ${failed} failed`);

    ElNotification({
        type: succeeded === results.length ? 'success' : 'warning',
        message: `Completed: ${succeeded}/${results.length} GIFs downloaded`,
        title: "Batch Download Complete"
    });

    return results;
}

/**
 * 获取队列状态
 */
export function getQueueStatus() {
    return workerManager.getStatus();
}

/**
 * 清空队列
 */
export function clearGIFQueue() {
    workerManager.clearQueue();
}

/**
 * 取消所有任务
 */
export function cancelAllGIFs() {
    workerManager.cancelAll();
}

// 导出管理器实例（用于高级用法）
export { workerManager };