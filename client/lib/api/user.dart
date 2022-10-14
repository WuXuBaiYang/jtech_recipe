import 'package:client/common/api/request.dart';
import 'package:client/common/common.dart';
import 'package:client/manage/auth.dart';
import 'package:client/model/model.dart';
import 'package:client/model/user.dart';
import 'package:client/tool/tool.dart';

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
      handle: (it) => UserModel.from(it),
    );
  }
}

// 单例调用
final userApi = UserAPI();
