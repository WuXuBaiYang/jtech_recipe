import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';
import 'base.dart';

/*
* 帖子标签接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class PostTagAPI extends BaseJAPI {
  // 添加标签
  Future<PostTagModel> add({
    required String text,
  }) {
    return handleResponseData(
      post("/post/tag",
          requestModel: RequestModel.body(
            data: {
              "name": text,
            },
          )),
      handle: (e) => PostTagModel.from(e),
    );
  }

  // 获取标签列表
  Future<PaginationModel<PostTagModel>> getList({
    required num pageIndex,
    int pageSize = 15,
    num userId = 0,
  }) {
    var path = "/post/tag";
    if (userId != 0) path = "$path/$userId";
    return handleResponseData(
        get(
          path,
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          ),
        ),
        handle: (e) => PaginationModel.from(
              e,
              parseItem: (it) => PostTagModel.from(it),
            ));
  }
}

// 单例调用
final postTagApi = PostTagAPI();
