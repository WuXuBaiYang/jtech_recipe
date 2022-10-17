import 'package:client/api/comment.dart';
import 'package:client/api/tag.dart';
import 'package:client/common/api/request.dart';
import 'package:client/model/comment.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';
import 'package:client/model/tag.dart';
import 'base.dart';

/*
* 帖子接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class PostAPI extends BaseJAPI {
  // 帖子发布
  Future<PostModel> createPost({
    required PostModel model,
  }) {
    return handleResponseData(
      post("/post",
          requestModel: RequestModel.body(
            data: model.toUpdateInfo(),
          )),
      handle: (e) => PostModel.from(e),
    );
  }

  // 编辑帖子
  Future<PostModel> updatePost({
    required PostModel model,
    required String postId,
  }) {
    return handleResponseData(
      put("/post/$postId",
          requestModel: RequestModel.body(
            data: model.toUpdateInfo(),
          )),
      handle: (e) => PostModel.from(e),
    );
  }

  // 获取帖子列表
  Future<PaginationModel<PostModel>> loadPosts({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    return handleResponsePaginationData(
      get("/post",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
              if (userId != null) "userId": userId
            },
          )),
      handle: (e) => PostModel.from(e),
    );
  }

  // 获取帖子详情
  Future<PostModel> loadPostInfo({
    required String postId,
  }) {
    return handleResponseData(
      get("/post/$postId"),
      handle: (e) => PostModel.from(e),
    );
  }

  // 发布帖子评论
  Future<CommentModel> createPostComment({
    required String postId,
    required String content,
  }) {
    return commentApi.createComment(
      path: "/post/comment",
      pId: postId,
      content: content,
    );
  }

  // 获取帖子评论列表
  Future<PaginationModel<CommentModel>> loadPostComments({
    required String postId,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return commentApi.loadComments(
      path: "/post/comment",
      pId: postId,
      pageIndex: pageIndex,
      pageSize: pageSize,
    );
  }

  // 批量添加帖子标签
  Future<List<TagModel>> addPostTags({
    required List<TagModel> tags,
  }) {
    return tagApi.addTags(
      path: "/post/tag",
      tags: tags,
    );
  }

  // 获取帖子标签集合
  Future<PaginationModel<TagModel>> loadPostTags({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    return tagApi.loadTags(
      path: "/post/tag",
      pageIndex: pageIndex,
      pageSize: pageSize,
      userId: userId,
    );
  }

  // 帖子点赞
  Future<bool> likePost({
    required String postId,
  }) {
    return handleResponseData(
      post("/post/like/$postId"),
    );
  }

  // 帖子取消点赞
  Future<bool> unLikePost({
    required String postId,
  }) {
    return handleResponseData(
      delete("/post/like/$postId"),
    );
  }

  // 帖子收藏
  Future<bool> collectPost({
    required String postId,
  }) {
    return handleResponseData(
      post("/post/collect/$postId"),
    );
  }

  // 帖子取消收藏
  Future<bool> unCollectPost({
    required String postId,
  }) {
    return handleResponseData(
      delete("/post/collect/$postId"),
    );
  }
}

// 单例调用
final postApi = PostAPI();
