import 'package:client/api/tag.dart';
import 'package:client/common/api/request.dart';
import 'package:client/model/menu.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';
import 'package:client/model/recipe.dart';
import 'package:client/model/tag.dart';
import 'package:client/model/user.dart';

import 'base.dart';

/*
* 用户接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class UserAPI extends BaseJAPI {
  // 获取用户信息
  Future<UserModel> loadUserInfo({
    String? userId,
  }) {
    var path = "/user/info";
    if (userId != null) path = "$path/$userId";
    return handleResponseData(
      get(path),
      handle: (e) => UserModel.from(e),
    );
  }

  // 更新用户信息
  Future<UserModel> updateUserInfo({
    required UserModel model,
  }) {
    return handleResponseData(
      put("/user/info",
          requestModel: RequestModel.body(
            data: model.toUpdateInfo(),
          )),
      handle: (e) => UserModel.from(e),
    );
  }

  // 批量添加收货地址标签
  Future<List<TagModel>> addUserAddressTags({
    required List<TagModel> tags,
  }) {
    return tagApi.addTags(
      path: "/user/tag/address",
      tags: tags,
    );
  }

  // 分页获取收货地址标签
  Future<PaginationModel<TagModel>> loadUserAddressTags({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return tagApi.loadTags(
      path: "/user/tag/address",
      pageIndex: pageIndex,
      pageSize: pageSize,
    );
  }

  // 添加收货地址
  Future<UserAddressModel> addUserAddress({
    required UserAddressModel model,
  }) {
    return handleResponseData(
      post(
        "/user/address",
        requestModel: RequestModel.body(
          data: model.toUpdateInfo(),
        ),
      ),
      handle: (e) => UserAddressModel.from(e),
    );
  }

  // 更新收货地址
  Future<UserAddressModel> updateUserAddress({
    required String addressId,
    required UserAddressModel model,
  }) {
    return handleResponseData(
      post(
        "/user/address/$addressId",
        requestModel: RequestModel.body(
          data: model.toUpdateInfo(),
        ),
      ),
      handle: (e) => UserAddressModel.from(e),
    );
  }

  // 修改收货地址为默认
  Future<bool> updateUserAddressDefault({
    required String addressId,
  }) {
    return handleResponseData(
      put("/user/address/$addressId/default"),
    );
  }

  // 修改收货地址排序
  Future<bool> updateUserAddressOrder({
    required String addressId,
    required int order,
  }) {
    return handleResponseData(
      put("/user/address/$addressId/order",
          requestModel: RequestModel.body(
            data: {
              "order": order,
            },
          )),
    );
  }

  // 获取全部收货地址
  Future<List<UserAddressModel>> loadAllUserAddress() {
    return handleResponseListData(
      get("/user/address"),
      handle: (e) => UserAddressModel.from(e),
    );
  }

  // 获取收货地址详情
  Future<UserAddressModel> loadUserAddressInfo({
    required String addressId,
  }) {
    return handleResponseData(
      get("/user/address/$addressId"),
      handle: (e) => UserAddressModel.from(e),
    );
  }

  // 获取我的帖子点赞列表
  Future<PaginationModel<PostModel>> loadLikePosts({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get("/user/post/like",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => PostModel.from(e),
    );
  }

  // 获取我的帖子收藏列表
  Future<PaginationModel<PostModel>> loadCollectPosts({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get("/user/post/collect",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => PostModel.from(e),
    );
  }

  // 获取我的菜单点赞列表
  Future<PaginationModel<MenuModel>> loadLikeMenus({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get("/user/menu/like",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 获取我的菜单收藏列表
  Future<PaginationModel<MenuModel>> loadCollectMenus({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get("/user/menu/collect",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => MenuModel.from(e),
    );
  }

  // 获取我的食谱点赞列表
  Future<PaginationModel<RecipeModel>> loadLikeRecipe({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get("/user/recipe/like",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => RecipeModel.from(e),
    );
  }

  // 获取我的食谱收藏列表
  Future<PaginationModel<RecipeModel>> loadCollectRecipe({
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return handleResponsePaginationData(
      get("/user/recipe/collect",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => RecipeModel.from(e),
    );
  }

  // 关注用户
  Future<bool> subscribeUser({
    required String userId,
  }) {
    return handleResponseData(
      post("/user/subscribe/$userId"),
    );
  }

  // 取消关注用户
  Future<bool> unSubscribeUser({
    required String userId,
  }) {
    return handleResponseData(
      delete("/user/subscribe/$userId"),
    );
  }

  // 获取用户关注列表
  Future<PaginationModel<UserModel>> loadSubscribeUsers({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    var path = "/user/subscribe";
    if (userId != null) path = "$path/$userId";
    return handleResponsePaginationData(
      get(path,
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
            },
          )),
      handle: (e) => UserModel.from(e),
    );
  }

  // 获取全部勋章
  Future<List<MedalModel>> loadAllMedals() {
    return handleResponseListData(
      get("/user/medal"),
      handle: (e) => MedalModel.from(e),
    );
  }
}

// 单例调用
final userApi = UserAPI();
