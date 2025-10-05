import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from "axios"
import { ElNotification } from "element-plus"

export const useNovelViewerStore = defineStore('novelViewer', () => {
    // 状态
    const isVisible = ref(false)
    const currentNovelId = ref('')
    const currentPage = ref(0)
    const totalPages = ref(1)
    const novelTitle = ref('')
    const author = ref('')
    const authorId = ref('')
    const isR18 = ref(false)
    const aiType = ref(false)
    const novelContent = ref('')
    const seriesData = ref(null)
    const tags = ref([])
    const description = ref('')
    const cover = ref('')

    const isLoading = ref(false)

    // 字体和样式设置
    const fontSize = ref(18)
    const minFontSize = ref(12)
    const maxFontSize = ref(32)
    const lineHeight = ref(1.8)

    // 计算属性
    const hasMultiplePages = computed(() => {
        return seriesData.value && (seriesData.value.prev || seriesData.value.next)
    })

    const canGoPrevious = computed(() => {
        return seriesData.value && seriesData.value.prev
    })

    const canGoNext = computed(() => {
        return seriesData.value && seriesData.value.next
    })

    // 方法
    const openViewer = async (novelData) => {
        console.log(novelData)
        currentNovelId.value = novelData.id
        currentPage.value = 0
        novelTitle.value = novelData.title || ''
        author.value = novelData.userName || ''
        authorId.value = novelData.userId || ''
        isR18.value = novelData.genre === '0' || false
        aiType.value = novelData.aiType || false
        seriesData.value = novelData.seriesNavData || null
        tags.value = novelData.tags || []
        description.value = novelData.description || ""
        cover.value = novelData.cover || ""

        // 如果有系列信息,计算总页数
        if (seriesData.value) {
            totalPages.value = calculateTotalPages(seriesData.value)
        } else {
            totalPages.value = 1
        }

        try {
            isVisible.value = true
            // 获取小说内容
            await loadNovelContent()

            // 禁用页面滚动
            document.body.style.overflow = 'hidden'
        } catch (error) {
            console.error('打开小说失败:', error)
            ElNotification({
                title: "小说打开失败",
                type: "warning",
                message: error.message || "请检查Cookie配置",
                position: 'bottom-right',
                duration: 3000,
            })
        }
    }

    const closeViewer = () => {
        isVisible.value = false
        currentNovelId.value = ''
        currentPage.value = 0
        totalPages.value = 1
        novelTitle.value = ''
        author.value = ''
        authorId.value = ''
        isR18.value = false
        aiType.value = false
        novelContent.value = ''
        seriesData.value = null
        tags.value = []
        description.value = ""
        cover.value = ""

        // 恢复页面滚动
        document.body.style.overflow = ''
    }

    const loadNovelContent = async () => {

        isLoading.value = true
        try {
            const response = await axios.get("http://127.0.0.1:7234/api/get_novel", {
                params: {
                    novelId: currentNovelId.value,
                }
            })
            let data =  JSON.parse(response.data.data)
            cover.value = data.coverUrl
            description.value = data.description
            novelTitle.value = data.title
            novelContent.value =data.content || data.body || ''
            seriesData.value = data.seriesNavData
            isR18.value = data.genre === '0' || false
            aiType.value = data.aiType || false
            let tmp = []
            data.tags.tags.forEach(it => {
                tmp.push(it.tag)
            })
            tags.value = tmp
            console.log(data)
        } catch (error) {
            console.error('加载小说内容失败:', error)
            throw new Error('加载内容失败,请稍后重试')
        } finally {
            isLoading.value = false
        }
    }

    const goToPrevious = async () => {
        if (canGoPrevious.value && seriesData.value.prev && seriesData.value.prev.available === true) {
            try {
                currentNovelId.value = seriesData.value.prev.id
                await loadNovelContent()
                currentPage.value -= 1
            } catch (error) {
                console.error('加载上一章失败:', error)
                ElNotification({
                    title: "加载失败",
                    type: "error",
                    message: "无法加载上一章",
                    position: 'bottom-right',
                    duration: 3000,
                })
            }
        }
    }

    const goToNext = async () => {
        if (canGoNext.value && seriesData.value.next && seriesData.value.next.available === true) {
            try {
                currentNovelId.value = seriesData.value.next.id
                await loadNovelContent()
                currentPage.value += 1
            } catch (error) {
                console.error('加载下一章失败:', error)
                ElNotification({
                    title: "加载失败",
                    type: "error",
                    message: "无法加载下一章",
                    position: 'bottom-right',
                    duration: 3000,
                })
            }
        }
    }

    const goToPage = async (page) => {
        if (page >= 0 && page < totalPages.value) {
            currentPage.value = page
            await loadNovelContent()
        }
    }

    // 字体控制
    const increaseFontSize = () => {
        if (fontSize.value < maxFontSize.value) {
            fontSize.value += 2
        }
    }

    const decreaseFontSize = () => {
        if (fontSize.value > minFontSize.value) {
            fontSize.value -= 2
        }
    }

    const setFontSize = (size) => {
        fontSize.value = Math.max(minFontSize.value, Math.min(size, maxFontSize.value))
    }

    const resetFontSize = () => {
        fontSize.value = 18
        lineHeight.value = 1.8
    }

    // 辅助函数:计算系列总页数
    const calculateTotalPages = (seriesData) => {
        // 这里简化处理,实际可能需要根据系列信息计算
        let count = 1
        if (seriesData.prev) count++
        if (seriesData.next) count++
        return count
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
            case '+':
            case '=':
                event.preventDefault()
                increaseFontSize()
                break
            case '-':
                event.preventDefault()
                decreaseFontSize()
                break
            case '0':
                event.preventDefault()
                resetFontSize()
                break
        }
    }

    return {
        // 状态
        isVisible,
        currentNovelId,
        currentPage,
        totalPages,
        novelTitle,
        author,
        authorId,
        isR18,
        aiType,
        novelContent,
        seriesData,
        tags,
        description,
        cover,
        isLoading,

        // 字体设置
        fontSize,
        minFontSize,
        maxFontSize,
        lineHeight,

        // 计算属性
        hasMultiplePages,
        canGoPrevious,
        canGoNext,

        // 方法
        openViewer,
        closeViewer,
        loadNovelContent,
        goToPrevious,
        goToNext,
        goToPage,
        handleKeyPress,

        // 字体方法
        increaseFontSize,
        decreaseFontSize,
        setFontSize,
        resetFontSize,
    }
})