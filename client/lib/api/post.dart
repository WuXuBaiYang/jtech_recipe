import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';

import 'base.dart';

/*
* 帖子接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class PostAPI extends BaseJAPI {
  // 帖子发布
  Future<PostModel> publish({
    required PostModel postInfo,
  }) {
    return handleResponseData(
      post(
        "/post",
        requestModel: RequestModel.body(
          data: postInfo.to(),
        ),
      ),
      handle: (e) => PostModel.from(e),
    );
  }

  // 帖子编辑
  Future<PostModel> update({
    required PostModel postInfo,
  }) {
    return handleResponseData(
      put(
        "/post",
        requestModel: RequestModel.body(
          data: postInfo.to(),
        ),
      ),
      handle: (e) => PostModel.from(e),
    );
  }

  // 帖子浏览
  Future<bool> view({
    required num postId,
  }) {
    return handleResponseData(
      post("/post/view/$postId"),
    );
  }

  // 帖子点赞
  Future<bool> like({
    required num postId,
  }) {
    return handleResponseData(
      post("/post/like/$postId"),
    );
  }

  // 帖子取消点赞
  Future<bool> unlike({
    required num postId,
  }) {
    return handleResponseData(
      delete("/post/like/$postId"),
    );
  }

  // 帖子收藏
  Future<bool> collect({
    required num postId,
  }) {
    return handleResponseData(
      post("/post/collect/$postId"),
    );
  }

  // 帖子取消收藏
  Future<bool> unCollect({
    required num postId,
  }) {
    return handleResponseData(
      delete("/post/collect/$postId"),
    );
  }

  // 获取帖子列表
  Future<PaginationModel<PostModel>> getList({
    required num pageIndex,
    int pageSize = 15,
    num userId = 0,
  }) {
    return handleResponseData(
      get(
        "post",
        requestModel: RequestModel.query(
          parameters: {
            "pageIndex": pageIndex,
            "pageSize": pageSize,
            if (userId != 0) "userId": userId,
          },
        ),
      ),
      handle: (e) => PaginationModel<PostModel>.from(
        e,
        parseItem: (it) => PostModel.from(it),
      ),
    );
  }

  // 获取帖子详情
  Future<PostModel> getInfo({
    required num postId,
  }) {
    return handleResponseData(
      get("/post/$postId"),
      handle: (e) => PostModel.from(e),
    );
  }
}

// 单例调用
final postApi = PostAPI();
