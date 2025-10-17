/**
 * 时间处理工具
 */

/**
 * 获取当前时间字符串
 * @param {boolean} withDate - 是否包含日期
 * @returns {string} 时间字符串
 */
export function getCurrentTime(withDate = false) {
    const now = new Date()
    if (withDate) {
        return now.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        })
    }
    return now.toLocaleTimeString('zh-CN')
}

/**
 * 格式化日期为 YYYY-MM-DD
 * @param {Date} date - 日期对象
 * @returns {string} 格式化的日期字符串
 */
export function formatDate(date) {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
}

/**
 * 获取今天的日期字符串
 * @returns {string} YYYY-MM-DD 格式
 */
export function getToday() {
    return formatDate(new Date())
}

/**
 * 获取 N 天前的日期
 * @param {number} days - 天数
 * @returns {string} YYYY-MM-DD 格式
 */
export function getDaysAgo(days) {
    const date = new Date()
    date.setDate(date.getDate() - days)
    return formatDate(date)
}

/**
 * 解析日期字符串
 * @param {string} dateStr - 日期字符串
 * @returns {Date} 日期对象
 */
export function parseDate(dateStr) {
    return new Date(dateStr)
}

/**
 * 计算两个日期之间的天数差
 * @param {Date|string} date1 - 日期1
 * @param {Date|string} date2 - 日期2
 * @returns {number} 天数差
 */
export function daysBetween(date1, date2) {
    const d1 = typeof date1 === 'string' ? parseDate(date1) : date1
    const d2 = typeof date2 === 'string' ? parseDate(date2) : date2
    const diff = Math.abs(d1 - d2)
    return Math.floor(diff / (1000 * 60 * 60 * 24))
}

/**
 * 格式化相对时间（如：2小时前）
 * @param {Date|string} date - 日期
 * @returns {string} 相对时间描述
 */
export function formatRelativeTime(date) {
    const now = new Date()
    const target = typeof date === 'string' ? parseDate(date) : date
    const diff = now - target

    const seconds = Math.floor(diff / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)

    if (days > 0) return `${days}天前`
    if (hours > 0) return `${hours}小时前`
    if (minutes > 0) return `${minutes}分钟前`
    return `${seconds}秒前`
}