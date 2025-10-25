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

/**
 * 解析日志行，提取类型、时间和内容
 * @param {string} logLine - 日志行文本，格式: [type] - HH:MM:SS message
 * @returns {Object|null} - 返回 {type, time, message} 或 null（解析失败）
 */
export function parseLogLine(logLine) {
    // 正则表达式: [类型] - 时间 内容
    const pattern = /^\[(\w+)\]\s*-\s*(\d{2}:\d{2}:\d{2})\s+(.*)$/;
    const match = logLine.match(pattern);

    if (!match) {
        return null; // 格式不匹配
    }

    return {
        type: match[1].toLowerCase(),  // 类型（转小写）
        time: match[2],                // 时间
        log: match[3].trim()       // 内容（去除首尾空格）
    };
}

// 示例3: 处理后端返回的日志数组
function processLogs(rawLogs) {
    return rawLogs.map(parseLogLine).filter(log => log !== null);
}

// 示例4: 按类型分组
function groupByType(logs) {
    const parsed = logs.map(parseLogLine).filter(log => log !== null);
    const grouped = {};

    parsed.forEach(log => {
        if (!grouped[log.type]) {
            grouped[log.type] = [];
        }
        grouped[log.type].push(log);
    });

    return grouped;
}