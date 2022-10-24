import 'dart:io';

import 'package:client/api/user.dart';
import 'package:client/common/common.dart';
import 'package:client/common/logic.dart';
import 'package:client/common/notifier.dart';
import 'package:client/manage/auth.dart';
import 'package:client/manage/oss.dart';
import 'package:client/manage/router.dart';
import 'package:client/model/user.dart';
import 'package:client/tool/snack.dart';
import 'package:client/widget/avatar.dart';
import 'package:client/widget/loading.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';

/*
* 授权初始化页
* @author wuxubaiyang
* @Time 2022/10/21 14:18
*/
class AuthInitPage extends StatefulWidget {
  const AuthInitPage({super.key});

  @override
  State<StatefulWidget> createState() => _AuthInitPageState();
}

/*
* 授权初始化页-状态页
* @author wuxubaiyang
* @Time 2022/10/21 14:18
*/
class _AuthInitPageState extends State<AuthInitPage> {
  // 逻辑管理
  final _logic = _AuthInitLogic();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: () =>
                    routerManage.pushReplacementNamed(RoutePath.home),
                child: const Text('跳过'),
              ),
            ),
            Form(
              key: _logic.formKey,
              child: Padding(
                padding: const EdgeInsets.all(45),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    FormField<String>(
                      initialValue: _logic.userInfo.avatar,
                      builder: (f) {
                        return Avatar.file(
                          file: File(f.value ?? ""),
                          avatarSize: AvatarSize.large,
                          onTap: () => _logic.pickAvatar(context).then(
                                (v) => f.didChange(v),
                              ),
                        );
                      },
                    ),
                    const SizedBox(height: 25),
                    TextFormField(
                      maxLength: 16,
                      onSaved: (v) => _logic.userInfo.nickName = v ?? '',
                      initialValue: _logic.userInfo.nickName,
                      decoration: const InputDecoration(
                        label: Text('昵称'),
                      ),
                    )
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
      floatingActionButton: ValueListenableBuilder<bool>(
        valueListenable: _logic.saveStateNotifier,
        builder: (_, authState, __) {
          return FloatingActionButton(
            onPressed: !authState
                ? () => _logic.saveUserInfo(context).then((v) {
                      if (v != null) {
                        routerManage.pushReplacementNamed(RoutePath.home);
                      }
                    })
                : null,
            child: Loading.dark(
              loading: authState,
              child: const Icon(Icons.done),
            ),
          );
        },
      ),
    );
  }
}

/*
* 授权初始化页-逻辑
* @author wuxubaiyang
* @Time 2022/10/21 14:25
*/
class _AuthInitLogic extends BaseLogic {
  // 表单key
  final formKey = GlobalKey<FormState>();

  // 用户信息
  final userInfo = authManage.userInfo;

  // 授权请求状态
  final saveStateNotifier = ValueChangeNotifier<bool>(false);

  // 保存用户信息
  Future<UserModel?> saveUserInfo(BuildContext context) async {
    final currState = formKey.currentState;
    if (currState == null || !currState.validate()) return null;
    try {
      currState.save();
      saveStateNotifier.setValue(true);
      // 提交头像图片到oss
      if (userInfo.avatar.isNotEmpty &&
          !userInfo.avatar.startsWith(OSSBucket.jTechRecipe.name)) {
        final avatar = await ossManage.uploadFile(
          File(userInfo.avatar),
        );
        if (avatar != null && avatar.isNotEmpty) {
          userInfo.avatar = avatar;
        }
      }
      return await userApi.updateUserInfo(model: userInfo);
    } catch (e) {
      SnackTool.showMessage(context, message: '保存用户信息失败');
      saveStateNotifier.setValue(false);
    }
    return null;
  }

  // 媒体选择器
  ImagePicker? _picker;

  // 头像选择
  Future<String> pickAvatar(BuildContext context) async {
    try {
      _picker ??= ImagePicker();
      final result = await _picker?.pickImage(
        source: ImageSource.gallery,
      );
      if (result != null) userInfo.avatar = result.path;
    } catch (e) {
      SnackTool.showMessage(context, message: '头像选择失败');
    }
    return userInfo.avatar;
  }
}
