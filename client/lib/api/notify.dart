import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/notify.dart';
import 'package:client/model/tag.dart';
import 'base.dart';

/*
* 通知接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class NotifyAPI extends BaseJAPI {
  // 获取通知集合
  Future<PaginationModel<NotifyModel>> loadNotifies({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get('/notification',
          requestModel: RequestModel.query(parameters: {
            'pageIndex': pageIndex,
            'pageSize': pageSize,
          })),
      handle: (e) => NotifyModel.from(e),
    );
  }
}

// 单例调用
final notifyApi = NotifyAPI();
