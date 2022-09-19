import 'package:client/api/base.dart';
import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/notification.dart';

/*
* 消息通知接口
* @author wuxubaiyang
* @Time 2022/9/13 15:26
*/
class NotificationApi extends BaseJAPI {
  // 分页获取系统通知
  Future<PaginationModel<NotificationModel>> getList({
    required num pageIndex,
    int pageSize = 15,
  }) {
    return handleResponseData(
      get(
        "/notification",
        requestModel: RequestModel.query(
          parameters: {
            "pageIndex": pageIndex,
            "pageSize": pageSize,
          },
        ),
      ),
      handle: (e) => PaginationModel<NotificationModel>.from(
        e,
        parseItem: (it) => NotificationModel.from(it),
      ),
    );
  }
}

// 单例调用
final notificationApi = NotificationApi();
