import 'dart:async';
import 'package:client/api/auth.dart';
import 'package:client/common/common.dart';
import 'package:client/common/logic.dart';
import 'package:client/common/notifier.dart';
import 'package:client/manage/router.dart';
import 'package:client/model/model.dart';
import 'package:client/tool/snack.dart';
import 'package:client/tool/tool.dart';
import 'package:client/widget/loading.dart';
import 'package:client/widget/listenable_builders.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/*
* 授权页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class AuthPage extends StatefulWidget {
  const AuthPage({super.key});

  @override
  State<StatefulWidget> createState() => _AuthPageState();
}

/*
* 授权页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _AuthPageState extends State<AuthPage> {
  // 逻辑管理
  final _logic = _AuthLogic();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        padding: const EdgeInsets.all(35),
        alignment: Alignment.center,
        child: _buildAuthForm(context),
      ),
      floatingActionButton: ValueListenableBuilder<bool>(
        valueListenable: _logic.authStateNotifier,
        builder: (_, authState, __) {
          return FloatingActionButton(
            onPressed: !authState
                ? () => _logic.authSaved(context).then((v) {
                      if (v != null) {
                        routerManage.pushReplacementNamed(
                            v.newUser ? RoutePath.authInit : RoutePath.home);
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

  // 构建授权表单
  Widget _buildAuthForm(BuildContext context) {
    return Form(
      key: _logic.formKey,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const FlutterLogo(size: 80),
          const SizedBox(height: 65),
          _buildPhoneField(),
          const SizedBox(height: 28),
          _buildSMSCodeField(context),
        ],
      ),
    );
  }

  // 构建手机号输入框
  Widget _buildPhoneField() {
    return TextFormField(
      autofocus: true,
      controller: _logic.phoneController,
      keyboardType: TextInputType.phone,
      textInputAction: TextInputAction.next,
      onChanged: (v) =>
          _logic.phoneVerifyNotifier.setValue(Tool.verifyPhone(v)),
      autovalidateMode: AutovalidateMode.onUserInteraction,
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
      ],
      decoration: InputDecoration(
        label: const Text('手机号'),
        hintText: '000 0000 0000',
        prefixIcon: const Icon(Icons.phone),
        suffixIcon: ValueListenableBuilder<bool>(
          valueListenable: _logic.phoneVerifyNotifier,
          builder: (_, v, __) {
            return Visibility(
              visible: v,
              child: const Icon(Icons.verified_outlined),
            );
          },
        ),
      ),
      validator: (v) {
        if (v == null || v.isEmpty) {
          return '手机号不能为空';
        }
        if (!Tool.verifyPhone(v)) {
          return '手机号校验失败';
        }
        return null;
      },
    );
  }

  // 构建短信验证码输入框
  Widget _buildSMSCodeField(BuildContext context) {
    return TextFormField(
      maxLength: 4,
      controller: _logic.smsCodeController,
      keyboardType: TextInputType.number,
      textInputAction: TextInputAction.done,
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
      ],
      autovalidateMode: AutovalidateMode.onUserInteraction,
      decoration: InputDecoration(
        label: const Text('验证码'),
        counterText: '',
        hintText: '0000',
        prefixIcon: const Icon(Icons.sms),
        suffixIcon: ValueListenableBuilder3<bool, int, SMSCodeState>(
          first: _logic.phoneVerifyNotifier,
          second: _logic.countdownSecondsNotifier,
          third: _logic.smsCodeStateNotifier,
          builder: (_, phoneVerified, countdownSeconds, smsCodeState, __) {
            final text = smsCodeState == SMSCodeState.loaded
                ? '验证码已发送($countdownSeconds)'
                : '获取验证码';
            return TextButton(
              onPressed: phoneVerified && smsCodeState == SMSCodeState.normal
                  ? () => _logic.sendSMS(context)
                  : null,
              child: Loading(
                loading: smsCodeState == SMSCodeState.loading,
                child: Text(text),
              ),
            );
          },
        ),
      ),
      validator: (v) {
        if (v == null || v.isEmpty) {
          return '验证码不能为空';
        }
        return null;
      },
    );
  }

  @override
  void dispose() {
    _logic.dispose();
    super.dispose();
  }
}

/*
* 授权页-逻辑
* @author wuxubaiyang
* @Time 2022/10/20 10:00
*/
class _AuthLogic extends BaseLogic {
  // 表单key
  final formKey = GlobalKey<FormState>();

  // 手机号输入框控制器
  final phoneController = TextEditingController();

  // 验证码输入框控制器
  final smsCodeController = TextEditingController();

  // 手机号校验状态
  final phoneVerifyNotifier = ValueChangeNotifier<bool>(false);

  // 授权请求状态
  final authStateNotifier = ValueChangeNotifier<bool>(false);

  // 授权信息保存
  Future<AuthModel?> authSaved(BuildContext context) async {
    final currState = formKey.currentState;
    if (currState == null || !currState.validate()) return null;
    try {
      authStateNotifier.setValue(true);
      final phoneNumber = phoneController.text;
      final code = smsCodeController.text;
      return await authApi.auth(phoneNumber: phoneNumber, code: code);
    } catch (e) {
      SnackTool.showMessage(context,
          message: kDebugMode ? e.toString() : '请求授权失败');
      authStateNotifier.setValue(false);
    }
    return null;
  }

  // 短信验证码状态
  final smsCodeStateNotifier =
      ValueChangeNotifier<SMSCodeState>(SMSCodeState.normal);

  // 发送短信验证码
  Future<void> sendSMS(BuildContext context) async {
    try {
      smsCodeStateNotifier.setValue(SMSCodeState.loading);
      final phoneNumber = phoneController.value.text;
      if (await authApi.sendSMS(phoneNumber: phoneNumber)) {
        _startSmsCountdown();
        if (kDebugMode) {
          // 开发模式，不会发送验证码，默认使用手机号后四位
          smsCodeController.text =
              phoneNumber.substring(phoneNumber.length - 4);
        }
      }
    } catch (e) {
      SnackTool.showMessage(context, message: '验证码发送失败');
      smsCodeStateNotifier.setValue(SMSCodeState.normal);
    }
  }

  // 记录短信验证码获取倒计时
  final countdownSecondsNotifier = ValueChangeNotifier<int>(0);

  // 计时器
  Timer? _countdownTimer;

  // 短信验证码获取倒计时
  void _startSmsCountdown() {
    smsCodeStateNotifier.setValue(SMSCodeState.loaded);
    final countDown = kDebugMode ? 5 : 60;
    countdownSecondsNotifier.setValue(countDown);
    _countdownTimer = Timer.periodic(
      const Duration(seconds: 1),
      (t) {
        final v = countdownSecondsNotifier.value;
        countdownSecondsNotifier.setValue(v - 1);
        if (countdownSecondsNotifier.value <= 0) {
          smsCodeStateNotifier.setValue(SMSCodeState.normal);
          t.cancel();
        }
      },
    );
  }

  @override
  void dispose() {
    // 销毁控制器和计时器
    _countdownTimer?.cancel();
    phoneController.dispose();
    smsCodeController.dispose();
    phoneVerifyNotifier.dispose();
    countdownSecondsNotifier.dispose();
    super.dispose();
  }
}

// 验证码提示按钮状态枚举
enum SMSCodeState { normal, loading, loaded }
