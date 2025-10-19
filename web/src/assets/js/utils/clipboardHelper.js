/**
 * 剪贴板操作工具
 */

import {ElNotification} from 'element-plus'

/**
 * 复制文本到剪贴板
 * @param {string} text - 要复制的文本
 * @param {string} label - 内容标签（用于通知）
 * @returns {Promise<boolean>} 是否成功
 */
export async function copyToClipboard(text, label = '内容') {
    try {
        await navigator.clipboard.writeText(text)
        ElNotification({
            position: "bottom-right",
            type: "success",
            message: `${label}已复制: ${text}`,
            duration: 2000,
        })
        return true
    } catch (err) {
        console.error('复制失败', err)
        ElNotification({
            position: "bottom-right",
            type: "warning",
            message: `${label}复制失败`,
            duration: 1000,
        })
        return false
    }
}

/**
 * 读取剪贴板内容
 * @returns {Promise<string>} 剪贴板文本
 */
export async function readFromClipboard() {
    try {
        const text = await navigator.clipboard.readText()
        return text
    } catch (err) {
        console.error('读取剪贴板失败', err)
        return ''
    }
}

/**
 * 复制链接到剪贴板并提示
 * @param {string} url - 链接
 * @returns {Promise<boolean>} 是否成功
 */
export async function copyLink(url) {
    return await copyToClipboard(url, '链接')
}