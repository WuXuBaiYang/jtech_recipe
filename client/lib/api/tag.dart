import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/tag.dart';
import 'package:client/model/user.dart';

import 'base.dart';

/*
* 用户接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class TagAPI extends BaseJAPI {
  // 批量添加标签
  Future<List<TagModel>> addTags({
    required String path,
    required List<TagModel> tags,
  }) {
    return handleResponseData(
      post(path,
          requestModel: RequestModel.body(
            data: {
              "dictList": tags.map((e) => e.toAddInfo()).toList(),
            },
          )),
      handle: (e) => (e ?? [])
          .map<TagModel>(
            (e) => TagModel.from(e),
          )
          .toList(),
    );
  }

  // 分页获取标签
  Future<PaginationModel<TagModel>> loadTags({
    required String path,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponseData(
        get(path,
            requestModel: RequestModel.query(parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            })),
        handle: (e) => PaginationModel.from(
              e,
              itemParse: (it) => TagModel.from(it),
            ));
  }
}

// 单例调用
final tagApi = TagAPI();
