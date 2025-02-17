API 文档

1. 通用说明

1.1 认证方式
所有需要认证的接口都需要在请求头中携带 token：
Authorization: Bearer <token>

1.2 响应格式
{
    "code": int,       // 状态码，200表示成功
    "msg": string,     // 状态信息
    "data": any        // 响应数据
}

1.3 分页参数
支持分页的接口都使用以下查询参数：
- page: int          // 页码，从1开始
- page_size: int     // 每页数量，最大500

2. API 接口

2.1 认证相关接口 (/api/auth)

2.1.1 登录
路径: POST /api/auth/login
请求体:
{
    "email": string,     // 邮箱，必填
    "password": string   // 密码，必填
}
响应:
{
    "code": 200,
    "msg": "ok",
    "data": {
        "token": string  // JWT token
    }
}

2.1.2 注册新用户 (需要管理员权限)
路径: POST /api/auth/register
请求体:
{
    "email": string,     // 邮箱，必填
    "password": string,  // 密码，可选，默认为 "123456"
    "name": string,      // 用户名，必填
    "role": string      // 角色，可选，admin或guest，默认为 "guest"
}
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.1.3 重置密码 (需要管理员权限)
路径: POST /api/auth/reset-password/:id
参数: 
- id: string           // 用户ID
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.1.4 删除用户 (需要管理员权限)
路径: DELETE /api/auth/:id
参数:
- id: string           // 用户ID
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.1.5 更新个人信息
路径: PUT /api/auth/self
请求体:
{
    "email": string,     // 邮箱，可选
    "password": string,  // 新密码，可选
    "name": string      // 新用户名，可选
}
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.1.6 获取所有用户列表 (需要管理员权限)
路径: GET /api/auth/list
响应:
{
    "code": 200,
    "msg": "ok",
    "data": [
        {
            "id": string,
            "name": string,
            "email": string,
            "role": string,
            "created_at": int64,
            "updated_at": int64,
            "created_by": string,
            "updated_by": string
        }
    ]
}

2.1.7 获取当前用户信息
路径: GET /api/auth/self
响应:
{
    "code": 200,
    "msg": "ok",
    "data": {
        "id": string,
        "name": string,
        "email": string,
        "role": string
    }
}

2.2 直播配置相关接口 (/api/live-conf)

2.2.1 创建配置 (需要管理员权限)
路径: POST /api/live-conf
请求体:
{
    "room_display_id": string,  // 房间显示ID，必填
    "url": string,             // 直播URL，必填
    "name": string,            // 配置名称，必填
    "enable": bool            // 是否启用，必填
}
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.2.2 更新配置 (需要管理员权限)
路径: PUT /api/live-conf
请求体:
{
    "id": int64,              // 配置ID，必填
    "room_display_id": string, // 房间显示ID，必填
    "url": string,            // 直播URL，必填
    "name": string,           // 配置名称，必填
    "enable": bool           // 是否启用，必填
}
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.2.3 删除配置 (需要管理员权限)
路径: DELETE /api/live-conf/:id
参数:
- id: int64            // 配置ID
响应:
{
    "code": 200,
    "msg": "ok",
    "data": null
}

2.2.4 获取单个配置
路径: GET /api/live-conf/:id
参数:
- id: int64            // 配置ID
响应:
{
    "code": 200,
    "msg": "ok",
    "data": {
        "id": int64,
        "room_display_id": string,
        "url": string,
        "name": string,
        "enable": bool,
        "modified_on": int64,
        "created_on": int64,
        "modified_by": string,
        "created_by": string
    }
}

2.2.5 获取配置列表
路径: GET /api/live-conf
响应:
{
    "code": 200,
    "msg": "ok",
    "data": {
        "list": [
            {
                "id": int64,
                "room_display_id": string,
                "url": string,
                "name": string,
                "enable": bool,
                "modified_on": int64,
                "created_on": int64,
                "modified_by": string,
                "created_by": string
            }
        ]
    }
}

2.3 礼物消息相关接口 (/api/gift-message)

2.3.1 获取礼物排行
路径: GET /api/gift-message/ranking
查询参数:
- to_user_ids: []uint64    // 接收用户ID列表，可选
- room_display_id: string  // 房间显示ID，必填
- begin: int64            // 开始时间戳，必填
- end: int64             // 结束时间戳，必填
响应:
{
    "code": 200,
    "msg": "ok",
    "data": [
        {
            "user_id": uint64,
            "user_name": string,
            "user_display_id": string,
            "total": int64,
            "room_display_id": string,
            "room_name": string,
            "to_user_id": uint64,
            "to_user_name": string,
            "to_user_display_id": string,
            "gift_list": [
                {
                    "gift_id": int64,
                    "gift_name": string,
                    "diamond_count": int64,
                    "combo_count": int64,
                    "image": string,
                    "message": string,
                    "timestamp": int64
                }
            ]
        }
    ]
}

2.3.2 获取接收用户列表
路径: GET /api/gift-message/to-user
查询参数:
- room_display_id: string  // 房间显示ID，必填
响应:
{
    "code": 200,
    "msg": "ok",
    "data": [
        {
            "id": uint64,
            "name": string,
            "display_id": string
        }
    ]
}

2.3.3 获取礼物消息列表
路径: GET /api/gift-message
查询参数:
- search: string          // 搜索关键词，可选
- user_ids: []uint64      // 发送用户ID列表，可选
- to_user_ids: []uint64   // 接收用户ID列表，可选
- room_display_id: string // 房间显示ID，必填
- begin: int64           // 开始时间戳，可选
- end: int64            // 结束时间戳，可选
- order_by: string       // 排序字段，可选
- order_direction: string // 排序方向(asc/desc)，可选
- diamond_count: int64   // 钻石数量，可选
- page: int             // 页码，必填，最小值1
- page_size: int        // 每页数量，必填，最小值1，最大值500
响应:
{
    "code": 200,
    "msg": "ok",
    "data": {
        "total": int64,
        "list": [
            {
                "id": int64,
                "user_id": uint64,
                "user_name": string,
                "user_display_id": string,
                "to_user_id": uint64,
                "to_user_name": string,
                "to_user_display_id": string,
                "gift_id": int64,
                "gift_name": string,
                "diamond_count": int64,
                "combo_count": string,
                "message": string,
                "timestamp": int64
            }
        ]
    }
}

2.4 普通消息相关接口 (/api/common-message)

2.4.1 获取消息列表
路径: GET /api/common-message
查询参数:
- room_display_id: string  // 房间显示ID，必填
- page: int              // 页码，必填，最小值1
- page_size: int         // 每页数量，必填，最小值1，最大值500
响应:
{
    "code": 200,
    "msg": "ok",
    "data": {
        "total": int64,
        "list": [
            {
                "id": int64,
                "user_id": uint64,
                "user_name": string,
                "user_display_id": string,
                "room_display_id": string,
                "message": string,
                "timestamp": int64
            }
        ]
    }
}

2.5 用户相关接口 (/api/user)

2.5.1 获取所有用户
路径: GET /api/user
响应:
{
    "code": 200,
    "msg": "ok",
    "data": [
        {
            "id": uint64,
            "name": string,
            "display_id": string
        }
    ]
}

2.5.2 搜索用户
路径: GET /api/user/search
查询参数:
- keyword: string        // 搜索关键词，必填
响应:
{
    "code": 200,
    "msg": "ok",
    "data": [
        {
            "id": uint64,
            "name": string,
            "display_id": string
        }
    ]
}
