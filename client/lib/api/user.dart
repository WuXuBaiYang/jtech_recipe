import 'package:client/common/api/request.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';
import 'package:client/model/user.dart';

import 'base.dart';

/*
* 用户接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class UserAPI extends BaseJAPI {
  // 获取用户信息
  Future<UserModel> getInfo({
    num userId = 0,
  }) {
    var path = "/user/info";
    if (userId != 0) path = "$path/$userId";
    return handleResponseData(
      get(path),
      handle: (it) => UserModel.from(it),
    );
  }

  // 更新用户信息
  Future<UserProfileModel> updateProfile({
    required UserProfileModel profile,
  }) {
    return handleResponseData(
      put(
        "/user/info",
        requestModel: RequestModel.body(
          data: profile.to(),
        ),
      ),
      handle: (e) => UserProfileModel.from(e),
    );
  }

  // 关注目标用户
  Future<bool> subscribe({
    required num userId,
  }) {
    return handleResponseData(
      post("/user/subscribe/$userId"),
    );
  }

  // 取消关注用户
  Future<bool> unsubscribe({
    required num userId,
  }) {
    return handleResponseData(
      delete("/user/subscribe/$userId"),
    );
  }

  // 获取已关注的用户列表
  Future<PaginationModel<UserModel>> getSubscribeList({
    required num pageIndex,
    int pageSize = 15,
    num userId = 0,
  }) {
    var path = "/user/subscribe";
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
      handle: (e) => PaginationModel<UserModel>.from(
        e,
        parseItem: (it) => UserModel.from(it),
      ),
    );
  }

  // 获取浏览帖子记录
  Future<PaginationModel<PostModel>> getPostViewList({
    required num pageIndex,
    int pageSize = 15,
    num userId = 0,
  }) {
    var path = "/user/common/view";
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
      handle: (e) => PaginationModel<PostModel>.from(
        e,
        parseItem: (it) => PostModel.from(it),
      ),
    );
  }

  // 获取点赞帖子记录
  Future<PaginationModel<PostModel>> getPostLikeList({
    required num pageIndex,
    int pageSize = 15,
    num userId = 0,
  }) {
    var path = "/user/common/like";
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
      handle: (e) => PaginationModel<PostModel>.from(
        e,
        parseItem: (it) => PostModel.from(it),
      ),
    );
  }

  // 获取收藏帖子记录
  Future<PaginationModel<PostModel>> getPostCollectList({
    required num pageIndex,
    int pageSize = 15,
    num userId = 0,
  }) {
    var path = "/user/common/collect";
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
      handle: (e) => PaginationModel<PostModel>.from(
        e,
        parseItem: (it) => PostModel.from(it),
      ),
    );
  }
}

// 单例调用
final userApi = UserAPI();
