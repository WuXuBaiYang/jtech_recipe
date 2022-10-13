import 'dart:async';
import 'dart:io';

import 'package:client/common/common.dart';
import 'package:client/common/manage.dart';
import 'package:client/manage/auth.dart';
import 'package:client/tool/log.dart';
import 'package:flutter_openim_sdk/flutter_openim_sdk.dart';

/*
* 即时通讯管理
* @author wuxubaiyang
* @Time 2022/9/9 17:28
*/
class IMManage extends BaseManage {
  static final IMManage _instance = IMManage._internal();

  factory IMManage() => _instance;

  IMManage._internal();

  @override
  Future<void> init() async {
    await OpenIM.iMManager.initSDK(
      platform: IMConfig.platformCode,
      apiAddr: IMConfig.apiAddr,
      wsAddr: IMConfig.wsAddr,
      dataDir: await IMConfig.dataDir ?? "",
      objectStorage: IMConfig.objectStorage,
      logLevel: IMConfig.logLevel,
      listener: IMConnectHandle(),
    );
    // 登录到im
    await loginIM();
  }

  // 判断是否已登录到聊天服务器
  bool get isLogin => OpenIM.iMManager.isLogined;

  // 使用授权信息登录到im
  Future<bool> loginIM() async {
    try {
      if (authManage.authorized) {
        var user = authManage.authInfo!.user;
        await OpenIM.iMManager.login(
          uid: user.imUserId,
          token: user.imToken,
          operationID: _genOperationID("login"),
        );
        return true;
      }
    } catch (e) {
      LogTool.e("im_login_error：", error: e);
    }
    return false;
  }

  // 退出登录
  Future<bool> logoutIM() async {
    try {
      await OpenIM.iMManager.logout(
        operationID: _genOperationID("logout"),
      );
      return true;
    } catch (e) {
      LogTool.e("im_logout_error：", error: e);
    }
    return false;
  }

  // 生成操作id
  String _genOperationID(String action) {
    var timestamp = DateTime.now().millisecondsSinceEpoch;
    return "${action}_${Platform.operatingSystem}_${authManage.userId}_$timestamp";
  }
}

// 单例调用
final imManage = IMManage();

/*
* 即时通信服务连接处理
* @author wuxubaiyang
* @Time 2022/9/9 18:02
*/
class IMConnectHandle extends OnConnectListener {
  // SDK连接服务器失败
  @override
  void connectFailed(int? code, String? errorMsg) {}

  // SDK连接服务器成功
  @override
  void connectSuccess() {}

  // SDK正在连接服务器
  @override
  void connecting() {}

  // 账号已在其他地方登录，当前设备被踢下线
  @override
  void kickedOffline() {}

  //  登录凭证过期，需要重新登录
  @override
  void userSigExpired() {}
}
