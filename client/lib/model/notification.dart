import 'package:client/common/model.dart';
import 'package:client/model/model.dart';

/*
* 消息通知结构体
* @author wuxubaiyang
* @Time 2022/9/12 20:08
*/
class NotificationModel extends BaseModel with BaseInfo, CreatorInfo {
  // 推送目标用户id
  final num targetUserId;

  // 消息类型（0系统消息|1关注对象发帖）
  final NotificationType type;

  // 消息标题
  final String title;

  // 消息内容
  final String content;

  // 消息路由
  final String uri;

  NotificationModel.from(obj)
      : targetUserId = obj?["targetUserId"] ?? 0,
        type = NotificationType.values[obj?["type"] ?? 0],
        title = obj?["title"] ?? "",
        content = obj?["content"] ?? "",
        uri = obj?["uri"] ?? "" {
    initialBaseInfo(obj);
    initialCreatorInfo(obj);
  }
}

/*
* 消息类型枚举
* @author wuxubaiyang
* @Time 2022/9/12 20:20
*/
enum NotificationType {
  // 系统消息
  system,
  // 关注对象发帖
  subscribePost
}

/*
* 消息类型枚举扩展
* @author wuxubaiyang
* @Time 2022/9/12 20:40
*/
extension NotificationTypeExtension on NotificationType {
  // 获取消息类型标题
  String get title {
    switch (this) {
      case NotificationType.system:
        return "系统消息";
      case NotificationType.subscribePost:
        return "关注对象发帖";
    }
  }
}
