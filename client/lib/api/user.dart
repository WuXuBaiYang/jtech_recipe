import 'package:client/api/tag.dart';
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
            data: model.toUserUpdateInfo(),
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
}

// 单例调用
final userApi = UserAPI();
