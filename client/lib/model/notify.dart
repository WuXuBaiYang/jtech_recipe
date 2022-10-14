import 'package:client/common/model.dart';
import 'package:client/model/model.dart';
import 'package:client/model/user.dart';

/*
* 消息通知
* @author wuxubaiyang
* @Time 2022/9/12 20:08
*/
class NotifyModel extends BaseModel with BasePart, CreatorPart {
  // 消息来源用户
  final String fromUserId;

  // 来源用户信息
  late UserModel? fromUser;

  // 目标用户id
  final String toUserId;

  // 消息类型
  final String typeCode;

  // 消息标题
  final String title;

  // 消息内容
  final String content;

  // 消息路由
  final String uri;

  NotifyModel.from(obj)
      : fromUserId = obj?["fromUserId"] ?? "",
        toUserId = obj?["toUserId"] ?? "",
        typeCode = obj?["typeCode"] ?? "",
        title = obj?["title"] ?? "",
        content = obj?["content"] ?? "",
        uri = obj?["uri"] ?? "" {
    initBasePart(obj);
    initCreatorPart(obj);
    if (obj?["fromUser"] != null) {
      fromUser = UserModel.from(obj?["fromUser"] ?? {});
    }
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        "fromUserId": fromUserId,
        "fromUser": fromUser?.to(),
        "toUserId": toUserId,
        "typeCode": typeCode,
        "title": title,
        "content": content,
        "uri": uri,
      };
}
