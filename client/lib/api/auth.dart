import 'package:client/common/api/request.dart';
import 'package:client/common/common.dart';
import 'package:client/manage/auth.dart';
import 'package:client/manage/im.dart';
import 'package:client/model/model.dart';
import 'package:client/tool/tool.dart';

import 'base.dart';

/*
* 授权接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class AuthAPI extends BaseJAPI {
  // 注册
  Future<AuthModel> register({
    required String userName,
    required String password,
  }) {
    return handleResponseData(
      post("/register",
          requestModel: RequestModel.body(data: {
            "userName": userName,
            "password": _signPassword(userName, password),
          })),
      handle: (it) => AuthModel.from(it),
    ).then(
      (v) => Future.wait([
        // 设置授权信息
        authManage.setupAuthInfo(v),
        // 登录到im
        imManage.loginIM(),
      ]).then((_) => v),
    );
  }

  // 登录
  Future<AuthModel> login({
    required String userName,
    required String password,
  }) {
    return handleResponseData(
      post("/login",
          requestModel: RequestModel.body(data: {
            "userName": userName,
            "password": _signPassword(userName, password),
          })),
      handle: (it) => AuthModel.from(it),
    ).then(
      (v) => Future.wait([
        // 设置授权信息
        authManage.setupAuthInfo(v),
        // 登录到im
        imManage.loginIM(),
      ]).then((_) => v),
    );
  }

  // 刷新token
  Future<AuthModel> refreshToken() {
    return handleResponseData(
      post("/refreshToken",
          requestModel: RequestModel.create(headers: {
            "RefreshToken": authManage.refreshToken,
          })),
      handle: (it) => AuthModel.from(it),
    ).then(
      (v) => Future.wait([
        // 设置授权信息
        authManage.setupAuthInfo(v),
      ]).then((_) => v),
    );
  }

  // 明文密码签名加密
  String _signPassword(String userName, String password) =>
      Tool.md5("$userName：${Common.salt}_${password}_${Common.salt}");
}

// 单例调用
final authApi = AuthAPI();
