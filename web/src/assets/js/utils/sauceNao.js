import { ElMessage } from 'element-plus'

export const SIMILARITY_CONFIG = {
    EXCELLENT: {
        min: 90,
        type: 'success',
        color: '#67C23A',
        icon: '🎯',
        text: '极高匹配',
        description: '几乎完全一致'
    },
    GOOD: {
        min: 75,
        type: 'primary',
        color: '#409EFF',
        icon: '✨',
        text: '良好匹配',
        description: '高度相似'
    },
    FAIR: {
        min: 60,
        type: 'warning',
        color: '#E6A23C',
        icon: '👌',
        text: '一般匹配',
        description: '可能相关'
    },
    POOR: {
        min: 0,
        type: 'info',
        color: '#909399',
        icon: '🤔',
        text: '较低匹配',
        description: '相似度低'
    }
}

export function getSimilarityConfig(similarity) {
    const sim = parseFloat(similarity)

    if (sim >= SIMILARITY_CONFIG.EXCELLENT.min) return SIMILARITY_CONFIG.EXCELLENT
    if (sim >= SIMILARITY_CONFIG.GOOD.min) return SIMILARITY_CONFIG.GOOD
    if (sim >= SIMILARITY_CONFIG.FAIR.min) return SIMILARITY_CONFIG.FAIR
    return SIMILARITY_CONFIG.POOR
}

export async function fetchSearch(file) {
    // 验证文件类型
    if (!file) {
        throw new Error('未选择文件')
    }

    const validTypes = ['image/jpeg', 'image/png']
    if (!validTypes.includes(file.type)) {
        throw new Error('只支持 JPG/PNG 格式的图片')
    }

    // 验证文件大小 (10MB)
    const maxSize = 10 * 1024 * 1024
    if (file.size > maxSize) {
        throw new Error('图片大小不能超过 10MB')
    }

    try {
        const formData = new FormData()
        formData.append('image', file)

        const response = await fetch('http://127.0.0.1:7235/api/saucenao/search', {
            method: 'POST',
            body: formData,
        })


        if (!response.ok) {
            const errorText = await response.text()
            throw new Error(`搜索请求失败: ${response.status} ${response.statusText} - ${errorText}`)
        }

        const data = await response.json()
        console.log(data)

        if (!data || !data.header || !data.results) {
            throw new Error('搜索结果格式错误')
        }

        return data

    } catch (error) {
        console.error('SauceNAO 搜索失败:', error)

        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            throw new Error('无法连接到服务器，请检查后端是否运行')
        } else if (error.message.includes('NetworkError')) {
            throw new Error('网络错误，请检查网络连接')
        } else {
            throw error
        }
    }
}


/**
 * 提取 Pixiv 作品 ID
 * @param {Object} result - SauceNAO 搜索结果项
 * @returns {string|null} Pixiv 作品 ID，如果不是 Pixiv 结果则返回 null
 */
export function extractPixivId(result) {
    console.log(result.data.pixiv_id)
    if (result.data.pixiv_id!=null && result.data.pixiv_id !== 0){
        return result.data.pixiv_id
    }

    const urls = []

    if (result.data.ext_urls && Array.isArray(result.data.ext_urls)) {
        urls.push(...result.data.ext_urls)
    }

    if (result.data.source) {
        if (typeof result.data.source === 'string') {
            urls.push(result.data.source)
        } else if (Array.isArray(result.data.source)) {
            urls.push(...result.data.source)
        }
    }

    for (const url of urls) {
        if (!url || typeof url !== 'string') continue

        // 匹配 pixiv.net/artworks/[id] 或 /illust/[id]
        const artworkMatch = url.match(/(?:artworks|illust)\/(\d+)/)
        if (artworkMatch) return artworkMatch[1]

        // 匹配 illust_id=[id] 参数
        const illustIdMatch = url.match(/illust_id=(\d+)/)
        if (illustIdMatch) return illustIdMatch[1]

        // 匹配 pximg.net 中的 ID（文件名格式: 123456_p0.jpg）
        const imgMatch = url.match(/\/(\d+)_p\d+\.(jpg|png|gif)/)
        if (imgMatch) return imgMatch[1]

        // 匹配 pximg.net 路径末尾的数字
        if (url.includes('pximg.net')) {
            const cleanUrl = url.split('?')[0].replace(/\.(jpg|png|gif|jpeg)$/i, '')
            const pathEndMatch = cleanUrl.match(/\/(\d{8,})$/)
            if (pathEndMatch) return pathEndMatch[1]
        }
    }

    return null
}

/**
 * 判断是否为 Pixiv 结果
 * @param {Object} result - SauceNAO 搜索结果项
 * @returns {boolean}
 */
export function isPixivResult(result) {
    if (!result.data.ext_urls ) return false
    let res = result.data.ext_urls.some(url =>
        url.includes('pixiv.net') || url.includes('pximg.net')
    )
    if (result.data.source != null){
        res = result.data.source.includes('pixiv.net') || result.data.source.includes('pximg.net')
    }
    return  res
}

/**
 * 获取搜索结果的标题
 * @param {Object} result - SauceNAO 搜索结果项
 * @returns {string}
 */
export function getResultTitle(result) {
    return result.data.title || result.data.source || '未知来源'
}

/**
 * 获取搜索结果的来源名称
 * @param {string} indexName - Index 名称，格式如 "Index #9: Danbooru - xxx.jpg"
 * @returns {string}
 */
export function getSourceName(indexName) {
    // 提取 "Index #X: Source Name" 中的 Source Name
    const match = indexName.match(/Index #\d+: (.+?)(?:\s-\s|$)/)
    return match ? match[1] : indexName
}

/**
 * 格式化 URL 以便显示
 * @param {string} url - 完整的 URL
 * @returns {string} 简化的域名
 */
export function formatUrl(url) {
    try {
        const urlObj = new URL(url)
        return urlObj.hostname.replace('www.', '')
    } catch {
        return url
    }
}

export default {
    fetchSearch,
    extractPixivId,
    isPixivResult,
    getResultTitle,
    getSourceName,
    getSimilarityConfig,
    formatUrl
}