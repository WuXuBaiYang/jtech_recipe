import 'package:client/api/comment.dart';
import 'package:client/api/tag.dart';
import 'package:client/common/api/request.dart';
import 'package:client/model/comment.dart';
import 'package:client/model/menu.dart';
import 'package:client/model/model.dart';
import 'package:client/model/tag.dart';
import 'base.dart';

/*
* 菜单接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class MenuAPI extends BaseJAPI {
  // 创建菜单
  Future<MenuModel> createMenu({
    required MenuModel model,
  }) {
    return handleResponseData(
      post("/menu",
          requestModel: RequestModel.body(
            data: model.toModifyInfo(),
          )),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 创建分支菜单
  Future<MenuModel> createForkMenu({
    required String menuId,
  }) {
    return handleResponseData(
      post("/menu/fork/$menuId"),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 编辑菜单
  Future<MenuModel> updateMenu({
    required MenuModel model,
    required String menuId,
  }) {
    return handleResponseData(
      put("/menu/$menuId",
          requestModel: RequestModel.body(
            data: model.toModifyInfo(),
          )),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 获取菜单集合
  Future<PaginationModel<MenuModel>> loadMenus({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    return handleResponsePaginationData(
      get("/menu",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
              if (userId != null) "userId": userId,
            },
          )),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 获取菜单详情
  Future<MenuModel> loadMenuInfo({
    required String menuId,
  }) {
    return handleResponseData(
      get("/menu/$menuId"),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 发布菜单评论
  Future<CommentModel> createMenuComment({
    required String menuId,
    required String content,
  }) {
    return commentApi.createComment(
      path: "/menu/comment",
      pId: menuId,
      content: content,
    );
  }

  // 获取菜单评论列表
  Future<PaginationModel<CommentModel>> loadPostComments({
    required String menuId,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return commentApi.loadComments(
      path: "/menu/comment",
      pId: menuId,
      pageIndex: pageIndex,
      pageSize: pageSize,
    );
  }

  // 批量添加菜单标签
  Future<List<TagModel>> addMenuTags({
    required List<TagModel> tags,
  }) {
    return tagApi.addTags(
      path: "/menu/tag",
      tags: tags,
    );
  }

  // 获取菜单标签集合
  Future<PaginationModel<TagModel>> loadMenuTags({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    return tagApi.loadTags(
      path: "/menu/tag",
      pageIndex: pageIndex,
      pageSize: pageSize,
      userId: userId,
    );
  }

  // 菜单点赞
  Future<bool> likeMenu({
    required String menuId,
  }) {
    return handleResponseData(
      post("/menu/like/$menuId"),
    );
  }

  // 菜单取消点赞
  Future<bool> unLikeMenu({
    required String menuId,
  }) {
    return handleResponseData(
      delete("/menu/like/$menuId"),
    );
  }

  // 菜单收藏
  Future<bool> collectPost({
    required String menuId,
  }) {
    return handleResponseData(
      post("/menu/collect/$menuId"),
    );
  }

  // 菜单取消收藏
  Future<bool> unCollectMenu({
    required String menuId,
  }) {
    return handleResponseData(
      delete("/menu/collect/$menuId"),
    );
  }
}

// 单例调用
final menuApi = MenuAPI();
