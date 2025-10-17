/**
 * 日志处理工具
 */

/**
 * 判断日志类型
 * @param {string} log - 日志内容
 * @returns {string} 日志类型类名
 */
export function getLogType(log) {
    if (log.includes('错误') || log.includes('Error')) {
        return 'log-error'
    } else if (log.includes('警告') || log.includes('Warning')) {
        return 'log-warning'
    } else if (log.includes('完成') || log.includes('Success')) {
        return 'log-success'
    }
    return 'log-info'
}

/**
 * 格式化日志条目
 * @param {string} message - 日志消息
 * @param {string} time - 时间戳（可选）
 * @returns {Object} 日志对象
 */
export function createLogEntry(message, time = null) {
    return {
        log: message,
        time: time || new Date().toLocaleTimeString(),
        type: getLogType(message)
    }
}

/**
 * 日志级别枚举
 */
export const LogLevel = {
    INFO: 'info',
    SUCCESS: 'success',
    WARNING: 'warning',
    ERROR: 'error'
}

/**
 * 创建指定级别的日志
 * @param {string} message - 消息
 * @param {string} level - 日志级别
 * @returns {Object} 日志对象
 */
export function createLog(message, level = LogLevel.INFO) {
    const typeMap = {
        [LogLevel.INFO]: 'log-info',
        [LogLevel.SUCCESS]: 'log-success',
        [LogLevel.WARNING]: 'log-warning',
        [LogLevel.ERROR]: 'log-error'
    }

    return {
        log: message,
        time: new Date().toLocaleTimeString(),
        type: typeMap[level] || 'log-info'
    }
}