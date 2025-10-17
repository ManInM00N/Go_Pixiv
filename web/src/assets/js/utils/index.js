/**
 * 工具函数统一导出
 * 使用示例：
 * import { formatCount, getImageUrl, copyToClipboard } from '@/utils'
 */

// 图片相关
export {
    getProxiedImageUrl,
    getImageUrl,
    getProfileUrl,
    getThumbnailUrl,
    preloadImage,
    preloadImages
} from './imageHelper.js'

// 格式化
export {
    formatCount,
    formatFileSize,
    formatDuration,
    formatPercentage,
    truncateText,
    isTextLong
} from './formatHelper.js'

// 剪贴板
export {
    copyToClipboard,
    readFromClipboard,
    copyLink
} from './clipboardHelper.js'

// 防抖节流
export {
    debounce,
    throttle,
    createDebouncedSearch,
    sleep
} from './debounce.js'

// 日志
export {
    getLogType,
    createLogEntry,
    createLog,
    LogLevel
} from './logHelper.js'

// 时间
export {
    getCurrentTime,
    formatDate,
    getToday,
    getDaysAgo,
    parseDate,
    daysBetween,
    formatRelativeTime
} from './timeHelper.js'

// Pixiv 相关
export {
    openPixivArtwork,
    openPixivUser,
    openPixivNovel,
    openPixivSeries,
    getPixivSearchUrl,
    isValidPixivId,
    extractPixivId,
    getRankingModeName
} from './pixivHelper.js'