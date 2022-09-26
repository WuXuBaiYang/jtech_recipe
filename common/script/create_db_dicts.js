mysql = require('mysql')
fs = require('fs')
flakeId = require('flake-idgen');
intformat = require('biguint-format')

flakeIdGen = new flakeId();

const tableMap = {
    'jtech_dict_index': './res/index_info.json',// 索引总表
    'jtech_dict_user_gender': './res/user_gender.json',// 用户性别
    'jtech_dict_user_evaluate': './res/user_evaluate.json',// 用户水平等级
    'jtech_dict_user_address_tag': './res/user_address_tag.json', // 用户收货地址标签
    'jtech_dict_notice_type': './res/notice_type.json',// 消息通知类型
    'jtech_dict_activity_type': './res/activity_type.json',// 活动类型
    'jtech_dict_post_tag': './res/post_tag.json',// 帖子标签
    'jtech_dict_address': './res/address.json',// 省市区三级联动
    'jtech_dict_recipe_cuisine': './res/recipe_cuisine.json',// 食谱菜系
    'jtech_dict_recipe_taste': './res/recipe_taste.json',// 食谱口味
    'jtech_dict_recipe_tag': './res/recipe_tag.json',// 食谱标签
    'jtech_dict_recipe_ingredients_main': './res/recipe_ingredients_main.json',// 食谱主材
    'jtech_dict_recipe_ingredients_sub': './res/recipe_ingredients_sub.json',// 食谱辅料
    'jtech_dict_medal_rarity': './res/medal_rarity.json',// 勋章稀有度等级
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
    let array = forEachResJson('', result)
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
        array.push(buildInsertItem(pCode, it.code, it.name, ""))
        if (it.children) array.push(...forEachResJson(it.code, it.children))
    }
    return array
}

// 构建表单插入map
function buildInsertItem(pCode, code, tag, desc) {
    let date = new Date()
    return [genId(), date, date, true, pCode, code, tag, desc]
}

// 批量插入数据
function batchInsert(name, values) {
    let sql = `insert into $name
                (id,created_at,updated_at,state,p_code,code,tag,$d) 
                values ?`.replaceAll('$name', name).replaceAll('$d', '\`desc\`')
    return new Promise((y, n) => {
        conn.query(sql, [values], function (err, result) {
            if (err) return n(new Error(name + '表数据插入失败：' + err.message))
            y(result)
        })
    })
}

// 创建表单
function createTable(name) {
    let sql = `create table if not exists $name
               (
                   id         bigint unsigned,
                   created_at datetime(3)     not null comment '创建时间',
                   updated_at datetime(3)     not null comment '更新时间',
                   deleted_at datetime(3)     null comment '删除时间',
                   creator_id bigint unsigned null comment '创建者id',
                   state      boolean         not null comment '是否可用',
                   p_code     varchar(20)     not null comment '父级字典码',
                   code       varchar(20)     not null unique comment '字典码',
                   tag        text            not null comment '标签',
                   $d       longtext        not null comment '描述',
                   primary key (id),
                   index idx_$name_deleted_at (deleted_at)
               )`.replaceAll('$name', name).replaceAll('$d', '\`desc\`')
    return new Promise((y, n) => {
        conn.query(sql, function (err, result) {
            if (err) return n(new Error(name + '表创建失败：' + err.message))
            y(result)
        })
    })
}

// 判断表是否存在
function tableExists(name) {
    let sql = `select *
               from information_schema.TABLES
               where TABLE_NAME = '$name'`
        .replaceAll('$name', name)
    return new Promise((y, n) => {
        conn.query(sql, function (err, result) {
            if (err) return n(new Error(name + '表查询失败：' + err.message))
            y(result.length > 0)
        })
    })
}

// 生成id
function genId() {
    return intformat(flakeIdGen.next(), 'dec')
}