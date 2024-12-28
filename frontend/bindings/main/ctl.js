// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

/**
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function CheckLogin() {
    let $resultPromise = /** @type {any} */($Call.ByID(980987062));
    return $resultPromise;
}

/**
 * @returns {Promise<void> & { cancel(): void }}
 */
export function Close() {
    let $resultPromise = /** @type {any} */($Call.ByID(2121317653));
    return $resultPromise;
}

/**
 * @param {string} text
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function DownloadByAuthorId(text) {
    let $resultPromise = /** @type {any} */($Call.ByID(3428520900, text));
    return $resultPromise;
}

/**
 * @param {string} page
 * @param {string} Type
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function DownloadByFollowPage(page, Type) {
    let $resultPromise = /** @type {any} */($Call.ByID(506309512, page, Type));
    return $resultPromise;
}

/**
 * @param {string} text
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function DownloadByNovelId(text) {
    let $resultPromise = /** @type {any} */($Call.ByID(1641432437, text));
    return $resultPromise;
}

/**
 * @param {string} text
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function DownloadByPid(text) {
    let $resultPromise = /** @type {any} */($Call.ByID(3069038031, text));
    return $resultPromise;
}

/**
 * @param {string} text
 * @param {string} Type
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function DownloadByRank(text, Type) {
    let $resultPromise = /** @type {any} */($Call.ByID(730038432, text, Type));
    return $resultPromise;
}

/**
 * @returns {Promise<string> & { cancel(): void }}
 */
export function ReturnString() {
    let $resultPromise = /** @type {any} */($Call.ByID(731407136));
    return $resultPromise;
}
