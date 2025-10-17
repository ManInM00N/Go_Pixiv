/**
 * 防抖和节流工具
 */

/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
export function debounce(func, delay = 300) {
    let timeoutId = null

    return function(...args) {
        clearTimeout(timeoutId)
        timeoutId = setTimeout(() => {
            func.apply(this, args)
        }, delay)
    }
}

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} 节流后的函数
 */
export function throttle(func, delay = 300) {
    let lastTime = 0

    return function(...args) {
        const now = Date.now()
        if (now - lastTime >= delay) {
            lastTime = now
            func.apply(this, args)
        }
    }
}

/**
 * 创建带有防抖的搜索处理器
 * @param {Function} searchFunc - 搜索函数
 * @param {number} delay - 延迟时间
 * @returns {Function} 防抖的搜索处理器
 */
export function createDebouncedSearch(searchFunc, delay = 300) {
    return debounce(searchFunc, delay)
}

/**
 * 延迟执行
 * @param {number} ms - 延迟毫秒数
 * @returns {Promise} Promise
 */
export function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
}