/**
 * 图片处理工具函数
 */

import { form } from '../configuration.js'
import noProfileImg from '../../images/NoR18.png'

/**
 * 获取代理后的图片 URL
 * @param {string} pixivUrl - Pixiv 原始图片 URL
 * @returns {string} 代理后的 URL
 */
export function getProxiedImageUrl(pixivUrl) {
    if (!pixivUrl) return ''
    const baseUrl = 'http://127.0.0.1:7234/api/preview?url='
    return `${baseUrl}${pixivUrl}`
}

/**
 * 获取图片 URL（带 R18 检查）
 * @param {string} imageUrl - 图片 URL
 * @param {number} isR18 - 是否为 R18 内容
 * @returns {string} 处理后的 URL
 */
export function getImageUrl(imageUrl, isR18 = 0) {
    if (isR18 && !form.value.r_18) {
        return noProfileImg
    }
    return getProxiedImageUrl(imageUrl)
}

/**
 * 获取头像 URL
 * @param {string} profileUrl - 头像 URL
 * @returns {string} 代理后的头像 URL
 */
export function getProfileUrl(profileUrl) {
    return getProxiedImageUrl(profileUrl)
}

/**
 * 获取缩略图 URL
 * @param {string} thumbUrl - 缩略图 URL
 * @returns {string} 代理后的缩略图 URL
 */
export function getThumbnailUrl(thumbUrl) {
    return getProxiedImageUrl(thumbUrl)
}

/**
 * 预加载图片
 * @param {string} url - 图片 URL
 * @returns {Promise} 加载完成的 Promise
 */
export function preloadImage(url) {
    return new Promise((resolve, reject) => {
        const img = new Image()
        img.onload = () => resolve(img)
        img.onerror = reject
        img.src = url
    })
}

/**
 * 批量预加载图片
 * @param {string[]} urls - 图片 URL 数组
 * @returns {Promise} 所有图片加载完成的 Promise
 */
export function preloadImages(urls) {
    return Promise.all(urls.map(url => preloadImage(url)))
}