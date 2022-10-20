import 'package:client/common/api/request.dart';
import 'package:client/model/comment.dart';
import 'package:client/model/model.dart';
import 'base.dart';

/*
* 评论/回复接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class CommentAPI extends BaseJAPI {
  // 发布评论
  Future<CommentModel> createComment({
    required String path,
    required String pId,
    required String content,
  }) {
    return handleResponseData(
      post(path,
          requestModel: RequestModel.body(
            data: {
              'content': content,
              'pId': pId,
            },
          )),
      handle: (e) => CommentModel.from(e),
    );
  }

  // 获取评论集合
  Future<PaginationModel<CommentModel>> loadComments({
    required String path,
    required String pId,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get(path,
          requestModel: RequestModel.query(
            parameters: {
              'pageIndex': pageIndex,
              'pageSize': pageSize,
              'pId': pId,
            },
          )),
      handle: (e) => CommentModel.from(e),
    );
  }

  // 发布回复
  Future<ReplayModel> createReplay({
    required ReplayModel model,
  }) {
    return handleResponseData(
      post('/replay',
          requestModel: RequestModel.body(
            data: model.toModifyInfo(),
          )),
      handle: (e) => ReplayModel.from(e),
    );
  }

  // 获取回复列表
  Future<PaginationModel<ReplayModel>> loadReplay({
    required String path,
    required String pId,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get(path,
          requestModel: RequestModel.query(
            parameters: {
              'pageIndex': pageIndex,
              'pageSize': pageSize,
              'pId': pId,
            },
          )),
      handle: (e) => ReplayModel.from(e),
    );
  }

  // 评论点赞
  Future<bool> likeComment({
    required String commentId,
  }) {
    return handleResponseData(
      post('/comment/like/$commentId'),
    );
  }

  // 评论取消点赞
  Future<bool> unLikeComment({
    required String commentId,
  }) {
    return handleResponseData(
      delete('/comment/like/$commentId'),
    );
  }

  // 回复点赞
  Future<bool> likeReplay({
    required String replayId,
  }) {
    return handleResponseData(
      post('/replay/like/$replayId'),
    );
  }

  // 回复取消点赞
  Future<bool> unLikeReplay({
    required String replayId,
  }) {
    return handleResponseData(
      delete('/replay/like/$replayId'),
    );
  }
}

// 单例调用
final commentApi = CommentAPI();
