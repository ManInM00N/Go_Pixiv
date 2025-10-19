/**
 * Pixiv 相关工具函数
 */

import {OpenInBrowser} from "../../../../bindings/main/internal/pixivlib/ctl.js";

/**
 * 打开 Pixiv 作品页面
 * @param {string|number} pid - 作品 ID
 */
export function openPixivArtwork(pid) {
    const url = `https://www.pixiv.net/artworks/${pid}`
    OpenInBrowser(url)
    window.open(url, '_blank')
}

/**
 * 打开 Pixiv 用户页面
 * @param {string|number} userId - 用户 ID
 */
export function openPixivUser(userId) {
    const url = `https://www.pixiv.net/users/${userId}`
    OpenInBrowser(url)
}

/**
 * 打开 Pixiv 小说页面
 * @param {string|number} novelId - 小说 ID
 */
export function openPixivNovel(novelId) {
    const url = `https://www.pixiv.net/novel/show.php?id=${novelId}`
    OpenInBrowser(url)
}

/**
 * 打开 Pixiv 系列页面
 * @param {string|number} seriesId - 系列 ID
 */
export function openPixivSeries(seriesId) {
    const url = `https://www.pixiv.net/novel/series/${seriesId}`
    OpenInBrowser(url)
}

/**
 * 生成 Pixiv 搜索 URL
 * @param {string} keyword - 搜索关键词
 * @param {string} type - 类型 (artworks, users, novels)
 * @returns {string} 搜索 URL
 */
export function getPixivSearchUrl(keyword, type = 'artworks') {
    const encodedKeyword = encodeURIComponent(keyword)
    return `https://www.pixiv.net/tags/${encodedKeyword}/${type}`

}

/**
 * 判断是否为有效的 Pixiv ID
 * @param {string|number} id - ID
 * @returns {boolean} 是否有效
 */
export function isValidPixivId(id) {
    const numId = Number(id)
    return !isNaN(numId) && numId > 0 && Number.isInteger(numId)
}

/**
 * 从 URL 提取 Pixiv ID
 * @param {string} url - Pixiv URL
 * @returns {string|null} 提取的 ID
 */
export function extractPixivId(url) {
    const patterns = [
        /artworks\/(\d+)/,
        /illust_id=(\d+)/,
        /novel\/show\.php\?id=(\d+)/,
        /users\/(\d+)/
    ]

    for (const pattern of patterns) {
        const match = url.match(pattern)
        if (match) return match[1]
    }

    return null
}

/**
 * 获取排行榜类型的中文名称
 * @param {string} mode - 排行榜模式
 * @returns {string} 中文名称
 */
export function getRankingModeName(mode) {
    const modeNames = {
        'daily': '每日排行',
        'weekly': '每周排行',
        'monthly': '每月排行',
        'daily_r18': '每日排行 R18',
        'weekly_r18': '每周排行 R18',
        'male': '男性向',
        'female': '女性向'
    }
    return modeNames[mode] || mode
}