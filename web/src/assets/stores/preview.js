import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from "axios";
import app from "../../App.vue";
import {ElNotification} from "element-plus";

export const useImageViewerStore = defineStore('imageViewer', () => {
    // 状态
    const isVisible = ref(false)
    const currentPid = ref('')
    const currentPage = ref(1)
    const totalPages = ref(1)
    const imageTitle = ref('')
    const thumbUrls = ref([])
    const regularUrls = ref([])
    const author = ref('')
    const authorId = ref('')
    const isR18 = ref(false)


    // 图片缩放和拖拽状态
    const scale = ref(1)
    const translateX = ref(0)
    const translateY = ref(0)
    const minScale = ref(0.1)
    const maxScale = ref(5)
    const isDragging = ref(false)
    const dragStartX = ref(0)
    const dragStartY = ref(0)
    const lastTranslateX = ref(0)
    const lastTranslateY = ref(0)

    // 计算属性
    const currentImageUrl = computed(() => {
        if (!currentPid.value) return ''
        return `${regularUrls.value[currentPage.value]}`
    })
    const hasMultiplePages = computed(() => totalPages.value > 1)

    const canGoPrevious = computed(() => currentPage.value > 0)

    const canGoNext = computed(() => currentPage.value + 1 < totalPages.value)

    const imageTransform = computed(() => {
        return `scale(${scale.value}) translate(${translateX.value}px, ${translateY.value}px)`
    })

    const isZoomedIn = computed(() => scale.value > 1)

    // 方法
    const  openViewer =async (imageData) => {
        currentPid.value = imageData.pid
        currentPage.value = 0
        totalPages.value = imageData.pages || 1
        imageTitle.value = imageData.title || ''
        author.value = imageData.author || ''
        authorId.value = imageData.authorId || ''
        isR18.value = imageData.r18 || false
        thumbUrls.value = []
        regularUrls.value = []
        resetTransform()
        await axios.get("http://127.0.0.1:7234/api/get_illust_page", {
            params: {
                pid: currentPid.value,
            }
        }).then((res) => {
            // console.log(res, res.data.body)
            let tmpUrl1 = [],tmpUrl2 = []
            for (var i = 0; i < res.data.body.length; i++) {
                tmpUrl1.push(res.data.body[i].urls.thumb_mini)
                tmpUrl2.push(res.data.body[i].urls.regular)
            }
            thumbUrls.value = tmpUrl1
            regularUrls.value = tmpUrl2
            console.log(thumbUrls.value,regularUrls.value)
            document.body.style.overflow = 'hidden'
        }).catch((error) => {
            console.log(error, error)
            ElNotification({
                title: "图片打开失败",
                type: "warning",
                message: error || "请检查Cookie配置",
                position: 'bottom-right',
                duration: 3000,
                }
            )
        }).finally(() => {

        })
        isVisible.value = true
    }

    const closeViewer = () => {
        isVisible.value = false
        currentPid.value = ''
        currentPage.value = 0
        totalPages.value = 1
        imageTitle.value = ''
        author.value = ''
        authorId.value = ''
        isR18.value = false

        resetTransform()

        // 恢复页面滚动
        document.body.style.overflow = ''
    }

    const goToPrevious = () => {
        if (canGoPrevious.value) {
            currentPage.value -= 1
        }
    }

    const goToNext = () => {
        if (canGoNext.value) {
            currentPage.value += 1
        }
    }

    const goToPage = (page) => {
        if (page >= 0 && page < totalPages.value) {
            currentPage.value = page
        }
    }


    const resetTransform = () => {
        scale.value = 1
        translateX.value = 0
        translateY.value = 0
        lastTranslateX.value = 0
        lastTranslateY.value = 0
    }

    const zoomIn = (step = 0.25) => {
        const newScale = Math.min(scale.value + step, maxScale.value)
        scale.value = newScale
    }

    const zoomOut = (step = 0.25) => {
        const newScale = Math.max(scale.value - step, minScale.value)
        scale.value = newScale

        // 如果缩小到1x以下，重置位移
        if (newScale <= 1) {
            translateX.value = 0
            translateY.value = 0
            lastTranslateX.value = 0
            lastTranslateY.value = 0
        }
    }

    const setZoom = (newScale) => {
        scale.value = Math.max(minScale.value, Math.min(newScale, maxScale.value))

        if (scale.value <= 1) {
            translateX.value = 0
            translateY.value = 0
            lastTranslateX.value = 0
            lastTranslateY.value = 0
        }
    }

    const fitToScreen = () => {
        setZoom(1)
    }

    // 拖拽开始
    const startDrag = (clientX, clientY) => {
        if (scale.value <= 1) return

        isDragging.value = true
        dragStartX.value = clientX
        dragStartY.value = clientY
        lastTranslateX.value = translateX.value
        lastTranslateY.value = translateY.value
    }

    // 拖拽移动
    const onDrag = (clientX, clientY) => {
        if (!isDragging.value || scale.value <= 1) return

        const deltaX = (clientX - dragStartX.value) / scale.value
        const deltaY = (clientY - dragStartY.value) / scale.value

        translateX.value = lastTranslateX.value + deltaX
        translateY.value = lastTranslateY.value + deltaY
    }

    // 拖拽结束
    const endDrag = () => {
        isDragging.value = false
    }

    // 滚轮缩放
    const handleWheel = (event, centerX = 0, centerY = 0) => {
        event.preventDefault()

        const delta = event.deltaY || event.detail || event.wheelDelta
        const zoomStep = 0.1
        const zoomDirection = delta > 0 ? -1 : 1

        const oldScale = scale.value
        const newScale = Math.max(minScale.value, Math.min(oldScale + (zoomStep * zoomDirection), maxScale.value))

        if (newScale !== oldScale) {
            // 计算以鼠标位置为中心的缩放
            const scaleChange = newScale / oldScale
            const rect = event.target.getBoundingClientRect()
            const offsetX = centerX - rect.width / 2
            const offsetY = centerY - rect.height / 2

            translateX.value = translateX.value * scaleChange - (offsetX * (scaleChange - 1)) / newScale
            translateY.value = translateY.value * scaleChange - (offsetY * (scaleChange - 1)) / newScale

            scale.value = newScale

            if (newScale <= 1) {
                translateX.value = 0
                translateY.value = 0
                lastTranslateX.value = 0
                lastTranslateY.value = 0
            }
        }
    }

    // 键盘快捷键处理
    const handleKeyPress = (event) => {
        if (!isVisible.value) return

        switch (event.key) {
            case 'Escape':
                closeViewer()
                break
            case 'ArrowLeft':
                event.preventDefault()
                goToPrevious()
                break
            case 'ArrowRight':
                event.preventDefault()
                goToNext()
                break
            case 'Home':
                event.preventDefault()
                goToPage(1)
                break
            case 'End':
                event.preventDefault()
                goToPage(totalPages.value)
                break
            case '+':
            case '=':
                event.preventDefault()
                zoomIn()
                break
            case '-':
                event.preventDefault()
                zoomOut()
                break
            case '0':
                event.preventDefault()
                fitToScreen()
                break
            case 'r':
                event.preventDefault()
                resetTransform()
                break
        }
    }

    return {
        // 状态
        isVisible,
        currentPid,
        currentPage,
        totalPages,
        imageTitle,
        author,
        authorId,
        isR18,
        thumbUrls,

        // 变换状态
        scale,
        translateX,
        translateY,
        isDragging,
        minScale,
        maxScale,

        // 计算属性
        currentImageUrl,
        hasMultiplePages,
        canGoPrevious,
        canGoNext,
        imageTransform,
        isZoomedIn,

        // 方法
        openViewer,
        closeViewer,
        goToPrevious,
        goToNext,
        goToPage,
        handleKeyPress,

        // 变换方法
        resetTransform,
        zoomIn,
        zoomOut,
        setZoom,
        fitToScreen,
        startDrag,
        onDrag,
        endDrag,
        handleWheel,
    }
})