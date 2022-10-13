import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';
import 'base.dart';

/*
* 帖子评论接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class PostCommentAPI extends BaseJAPI {
  // 发布评论
  Future<PostCommentModel> publish({
    required num postId,
    required String content,
  }) {
    return handleResponseData(
      post(
        "/post/comment/$postId",
        requestModel: RequestModel.body(
          data: {
            "content": content,
          },
        ),
      ),
      handle: (e) => PostCommentModel.from(e),
    );
  }

  // 获取评论列表
  Future<PaginationModel<PostCommentModel>> getList({
    required num postId,
    required num pageIndex,
    int pageSize = 15,
  }) {
    return handleResponseData(
      get(
        "/post/comment/$postId",
        requestModel: RequestModel.query(
          parameters: {
            "pageIndex": pageIndex,
            "pageSize": pageSize,
          },
        ),
      ),
      handle: (e) => PaginationModel<PostCommentModel>.from(
        e,
        parseItem: (it) => PostCommentModel.from(it),
      ),
    );
  }

  // 评论点赞
  Future<bool> like({
    required num commentId,
  }) {
    return handleResponseData(
      post("/post/comment/like/$commentId"),
    );
  }

  // 评论取消点赞
  Future<bool> unlike({
    required num commentId,
  }) {
    return handleResponseData(
      delete("/post/comment/like/$commentId"),
    );
  }
}

// 单例调用
final postCommentApi = PostCommentAPI();

/*
* 帖子评论回复接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class PostCommentReplayAPI extends BaseJAPI {
  // 发布评论回复
  Future<PostCommentReplayModel> publish({
    required num commentId,
    required String content,
  }) {
    return handleResponseData(
      post(
        "/post/comment/replay/$commentId",
        requestModel: RequestModel.body(
          data: {
            "content": content,
          },
        ),
      ),
      handle: (e) => PostCommentReplayModel.from(e),
    );
  }

  // 获取评论回复列表
  Future<PaginationModel<PostCommentReplayModel>> getList({
    required num commentId,
    required num pageIndex,
    int pageSize = 15,
  }) {
    return handleResponseData(
      get(
        "/post/comment/replay/$commentId",
        requestModel: RequestModel.query(
          parameters: {
            "pageIndex": pageIndex,
            "pageSize": pageSize,
          },
        ),
      ),
      handle: (e) => PaginationModel<PostCommentReplayModel>.from(
        e,
        parseItem: (it) => PostCommentReplayModel.from(it),
      ),
    );
  }

  // 评论回复点赞
  Future<bool> like({
    required num replayId,
  }) {
    return handleResponseData(
      post("/post/comment/replay/like/$replayId"),
    );
  }

  // 评论回复取消点赞
  Future<bool> unlike({
    required num replayId,
  }) {
    return handleResponseData(
      delete("/post/comment/replay/like/$replayId"),
    );
  }
}

// 单例调用
final postCommentReplayAPI = PostCommentReplayAPI();
