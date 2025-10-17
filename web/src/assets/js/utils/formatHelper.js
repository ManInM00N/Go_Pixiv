/**
 * 格式化工具函数
 */

/**
 * 格式化数字（1000 -> 1K, 10000 -> 1W）
 * @param {number} count - 要格式化的数字
 * @returns {string} 格式化后的字符串
 */
export function formatCount(count) {
    if (!count || count < 1000) return count?.toString() || '0'
    if (count < 10000) return (count / 1000).toFixed(1) + 'K'
    if (count < 100000) return (count / 10000).toFixed(1) + 'W'
    return (count / 10000).toFixed(0) + 'W'
}

/**
 * 格式化文件大小
 * @param {number} bytes - 字节数
 * @returns {string} 格式化后的大小字符串
 */
export function formatFileSize(bytes) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

/**
 * 格式化时长（秒数转为 HH:MM:SS）
 * @param {number} seconds - 秒数
 * @returns {string} 格式化后的时长字符串
 */
export function formatDuration(seconds) {
    const h = Math.floor(seconds / 3600)
    const m = Math.floor((seconds % 3600) / 60)
    const s = Math.floor(seconds % 60)

    const parts = []
    if (h > 0) parts.push(h.toString().padStart(2, '0'))
    parts.push(m.toString().padStart(2, '0'))
    parts.push(s.toString().padStart(2, '0'))

    return parts.join(':')
}

/**
 * 格式化百分比
 * @param {number} value - 数值
 * @param {number} total - 总数
 * @param {number} decimals - 小数位数
 * @returns {string} 百分比字符串
 */
export function formatPercentage(value, total, decimals = 1) {
    if (total === 0) return '0%'
    return ((value / total) * 100).toFixed(decimals) + '%'
}

/**
 * 截断文本
 * @param {string} text - 文本
 * @param {number} maxLength - 最大长度
 * @param {string} suffix - 后缀
 * @returns {string} 截断后的文本
 */
export function truncateText(text, maxLength = 50, suffix = '...') {
    if (!text || text.length <= maxLength) return text
    return text.substring(0, maxLength) + suffix
}

/**
 * 判断文本是否过长
 * @param {string} text - 文本
 * @param {number} maxLength - 最大长度
 * @returns {boolean} 是否过长
 */
export function isTextLong(text, maxLength = 30) {
    return text && text.length > maxLength
}