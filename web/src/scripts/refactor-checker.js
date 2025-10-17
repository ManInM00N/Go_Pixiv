#!/usr/bin/env node
/**
 * 简化版重构检查脚本 (ES Module)
 * 使用方法: node scripts/simple-checker.js
 */

import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

// 获取当前文件的目录
const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

// 颜色输出
const colors = {
    reset: '\x1b[0m',
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    magenta: '\x1b[35m'
}

function log(message, color = 'reset') {
    console.log(`${colors[color]}${message}${colors.reset}`)
}

// 配置
const paths = {
    utils: path.join(__dirname, '../assets/js/utils'),
    styles: path.join(__dirname, '../assets/style/common'),
    components: path.join(__dirname, '../components')
}

// 需要检查的文件
const requiredFiles = {
    utils: [
        'index.js',
        'imageHelper.js',
        'formatHelper.js',
        'clipboardHelper.js',
        'debounce.js',
        'logHelper.js',
        'timeHelper.js',
        'pixivHelper.js'
    ],
    styles: [
        'index.less',
        'page-header.less',
        'cards.less',
        'buttons.less',
        'animations.less',
        'modal.less',
        'loading.less',
        'waterfall.less',
        'pagination.less',
        'responsive.less'
    ]
}

// 检查文件是否存在
function checkFiles(type) {
    const dir = paths[type]
    const files = requiredFiles[type]

    log(`\n📁 检查 ${type === 'utils' ? '工具函数' : '样式模块'}...`, 'blue')

    if (!fs.existsSync(dir)) {
        log(`  ✗ 目录不存在: ${dir}`, 'red')
        return { passed: 0, failed: files.length }
    }

    let passed = 0
    let failed = 0

    files.forEach(file => {
        const filePath = path.join(dir, file)
        if (fs.existsSync(filePath)) {
            log(`  ✓ ${file}`, 'green')
            passed++
        } else {
            log(`  ✗ ${file} - 文件不存在`, 'red')
            failed++
        }
    })

    const total = files.length
    const percentage = ((passed / total) * 100).toFixed(0)
    const statusColor = passed === total ? 'green' : failed > passed ? 'red' : 'yellow'

    log(`  结果: ${passed}/${total} (${percentage}%) 通过`, statusColor)

    return { passed, failed, total }
}

// 检查组件是否导入了工具函数
function checkComponentImports() {
    log(`\n📝 检查组件导入...`, 'blue')

    const componentsToCheck = [
        'follow.vue',
        'PicCard.vue',
        'NovelCard.vue',
        'NovelMask.vue',
        'PicMask.vue',
        'maindownload.vue',
        'rank.vue',
        'NovelPage.vue',
        'settings.vue'
    ]

    let passed = 0
    let warnings = 0

    componentsToCheck.forEach(file => {
        const filePath = path.join(paths.components, file)

        if (!fs.existsSync(filePath)) {
            log(`  ⚠ ${file} - 文件不存在`, 'yellow')
            warnings++
            return
        }

        const content = fs.readFileSync(filePath, 'utf-8')

        // 检查是否有 utils 导入
        const hasUtilsImport = /import\s+.*from\s+['"].*utils['"]/g.test(content)
        // 检查是否有样式导入
        const hasStyleImport = /@import\s+["'].*\/common\//g.test(content)

        // Menu_List 和 search 可以跳过
        if (file === 'Menu_List.vue' || file === 'search.vue') {
            log(`  ○ ${file} - 跳过检查`, 'blue')
            passed++
            return
        }

        const issues = []
        if (!hasUtilsImport) {
            issues.push('可能缺少工具函数导入')
        }
        if (!hasStyleImport) {
            issues.push('可能缺少样式模块导入')
        }

        if (issues.length === 0) {
            log(`  ✓ ${file}`, 'green')
            passed++
        } else {
            log(`  ⚠ ${file} - ${issues.join(', ')}`, 'yellow')
            warnings++
        }
    })

    const total = componentsToCheck.length
    log(`  结果: ${passed} 通过, ${warnings} 警告`, warnings === 0 ? 'green' : 'yellow')

    return { passed, warnings, total }
}

// 生成报告
function generateReport(results) {
    log(`\n${'='.repeat(60)}`, 'magenta')
    log(`📊 重构检查报告`, 'magenta')
    log(`${'='.repeat(60)}`, 'magenta')

    const totalPassed = results.utils.passed + results.styles.passed + results.components.passed
    const totalIssues = results.utils.failed + results.styles.failed + results.components.warnings

    log(`\n总体情况:`, 'blue')
    log(`  ✓ 通过: ${totalPassed}`, 'green')

    if (totalIssues > 0) {
        log(`  ⚠ 问题: ${totalIssues}`, 'yellow')
    }

    log(`\n详细统计:`, 'blue')
    log(`  工具函数模块: ${results.utils.passed}/${results.utils.total}`)
    log(`  样式模块: ${results.styles.passed}/${results.styles.total}`)
    log(`  组件检查: ${results.components.passed}/${results.components.total}`)

    if (totalIssues === 0) {
        log(`\n🎉 所有检查通过！`, 'green')
    } else {
        log(`\n💡 发现 ${totalIssues} 个问题，请查看详细信息`, 'yellow')
    }

    log(`${'='.repeat(60)}\n`, 'magenta')
}

// 主函数
function main() {
    log(`\n🔍 开始检查重构状态...`, 'blue')
    log(`时间: ${new Date().toLocaleString()}\n`)

    const results = {
        utils: checkFiles('utils'),
        styles: checkFiles('styles'),
        components: checkComponentImports()
    }

    generateReport(results)

    // 返回退出码
    const hasCriticalIssues = results.utils.failed > 0 || results.styles.failed > 0
    process.exit(hasCriticalIssues ? 1 : 0)
}

// 运行
main()