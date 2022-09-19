mysql = require('mysql')
fs = require('fs')

const tableMap = {
    'dict_user_gender': './res/user_gender.json',// 用户性别
    'dict_user_evaluate': './res/user_evaluate.json',// 用户水平等级
    'dict_user_address_tag': './res/user_address_tag.json', // 用户收货地址标签
    'dict_notice_type': './res/notice_type.json',// 消息通知类型
    'dict_activity_type': './res/activity_type.json',// 活动类型
    'dict_post_tag': './res/post_tag.json',// 帖子标签
    'dict_address': './res/address.json',// 省市区三级联动
    'dict_recipe_cuisine': './res/recipe_cuisine.json',// 食谱菜系
    'dict_recipe_taste': './res/recipe_taste.json',// 食谱口味
    'dict_recipe_tag': './res/recipe_tag.json',// 食谱标签
    'dict_recipe_not_suitable': './res/recipe_not_suitable.json',// 食谱禁忌
    'dict_recipe_ingredients_main': './res/recipe_ingredients_main.json',// 食谱主材
    'dict_recipe_ingredients_sub': './res/recipe_ingredients_main.json',// 食谱辅料
    'dict_recipe_supplies': './res/recipe_supplies.json',// 食谱器具
}

conn = mysql.createConnection({
    host: 'localhost',// 也可以连接正式服务器
    user: 'jtech_server', password: 'JXuIAi4wqP0kho', database: 'jtech_recipe',
})
conn.connect(async function (err) {
    if (err) return console.error('数据库连接失败：' + err.message)
    console.log('数据库已连接')
    console.log('*数据表已存在时不会执行任何操作，如需重新创建请删除该表')
    try {
        for (const k in tableMap) {
            console.log('---------------  ' + k + '  ---------------')
            let v = tableMap[k]
            await handleDict(k, v)
        }
        console.log('---------------  所有操作已完成  ---------------')
    } catch (e) {
        console.log(e)
    }
    process.exit(0)
})

// 处理字典表
async function handleDict(name, resPath) {
    if (await tableExists(name)) return console.log('表已存在')
    console.log('正在创建字典表')
    await createTable(name)
    console.log('正在读取资源文件')
    let result = await fs.readFileSync(resPath)
    result = JSON.parse(result.toString())
    let array = forEachResJson(0, result)
    if (array.length < 1000) {
        console.log('正在插入数据（' + array.length + '）')
        await batchInsert(name, array)
    } else {
        console.log('拆分集合分批次插入')
        const subGroupLength = 1000
        let index = 0
        let i = 1
        while (index < array.length) {
            let newArray = array.slice(index, index += subGroupLength)
            console.log('正在插入第' + i++ + '组数据（' + newArray.length + '）')
            await batchInsert(name, newArray)
        }
    }
    console.log('数据插入成功')
}

// 迭代资源中的json文件
function forEachResJson(pCode, result) {
    let array = []
    for (const i in result) {
        let it = result[i]
        let code = parseInt(it.code)
        array.push(buildInsertItem(pCode, it.name, code, ""))
        if (it.children) array.push(...forEachResJson(code, it.children))
    }
    return array
}

// 构建表单插入map
function buildInsertItem(pCode, tag, code, desc) {
    let date = new Date()
    return [date, date, pCode, 0, tag, 0, code, true, desc]
}

// 批量插入数据
function batchInsert(name, values) {
    let sql = "INSERT INTO `" + name + "` (`created_at`,`updated_at`,`p_code`,`creator_id`,`tag`,`order`,`code`,`state`,`desc`) VALUES ?";
    return new Promise((y, n) => {
        conn.query(sql, [values], function (err, result) {
            if (err) return n(new Error(name + '表数据插入失败：' + err.message))
            y(result)
        })
    })
}

// 创建表单
function createTable(name) {
    let sql = "CREATE TABLE IF NOT EXISTS `" + name + "` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL COMMENT '创建时间',`updated_at` datetime(3) NULL COMMENT '更新时间',`deleted_at` datetime(3) NULL COMMENT '删除时间',`p_code` bigint unsigned NOT NULL COMMENT '父id',`creator_id` bigint unsigned COMMENT '创建者id',`tag` longtext NOT NULL COMMENT '标签',`order` bigint NOT NULL COMMENT '排序',`code` bigint NOT NULL UNIQUE COMMENT '键值',`state` boolean NOT NULL COMMENT '是否可用',`desc` longtext COMMENT '标记/描述',PRIMARY KEY (`id`),INDEX `idx_dict_deleted_at` (`deleted_at`))"
    return new Promise((y, n) => {
        conn.query(sql, function (err, result) {
            if (err) return n(new Error(name + '表创建失败：' + err.message))
            y(result)
        })
    })
}

// 判断表是否存在
function tableExists(name) {
    let sql = "select * from information_schema.TABLES where TABLE_NAME = '" + name + "'"
    return new Promise((y, n) => {
        conn.query(sql, function (err, result) {
            if (err) return n(new Error(name + '表查询失败：' + err.message))
            y(result.length > 0)
        })
    })
}