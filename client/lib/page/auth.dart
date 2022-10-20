import 'dart:async';
import 'package:client/api/auth.dart';
import 'package:client/common/common.dart';
import 'package:client/common/logic.dart';
import 'package:client/common/notifier.dart';
import 'package:client/main.dart';
import 'package:client/manage/router.dart';
import 'package:client/tool/snack.dart';
import 'package:client/tool/tool.dart';
import 'package:client/widget/loading.dart';
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
  // 页面逻辑管理
  final logic = _AuthLogic();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        padding: const EdgeInsets.all(35),
        alignment: Alignment.center,
        child: _buildAuthForm(context),
      ),
      floatingActionButton: ValueListenableBuilder<bool>(
        valueListenable: logic.authStateNotifier,
        builder: (_, authState, __) {
          return FloatingActionButton(
            onPressed: !authState ? () => logic.authSaved(context) : null,
            child: LoadingView.dark(
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
      key: logic.formKey,
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
      controller: logic.phoneController,
      keyboardType: TextInputType.phone,
      textInputAction: TextInputAction.next,
      onChanged: (v) => logic.phoneVerifyNotifier.setValue(Tool.verifyPhone(v)),
      autovalidateMode: AutovalidateMode.onUserInteraction,
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
      ],
      decoration: InputDecoration(
        label: const Text("手机号"),
        hintText: "000 0000 0000",
        prefixIcon: const Icon(Icons.phone),
        suffixIcon: ValueListenableBuilder<bool>(
          valueListenable: logic.phoneVerifyNotifier,
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
          return "手机号不能为空";
        }
        if (!Tool.verifyPhone(v)) {
          return "手机号校验失败";
        }
        return null;
      },
    );
  }

  // 构建短信验证码输入框
  Widget _buildSMSCodeField(BuildContext context) {
    return TextFormField(
      maxLength: 4,
      controller: logic.smsCodeController,
      keyboardType: TextInputType.number,
      textInputAction: TextInputAction.done,
      onFieldSubmitted: (v) => logic.authSaved(context),
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
      ],
      autovalidateMode: AutovalidateMode.onUserInteraction,
      decoration: InputDecoration(
        counter: const SizedBox(),
        label: const Text("验证码"),
        hintText: "0000",
        prefixIcon: const Icon(Icons.sms),
        suffixIcon: ValueListenableBuilder<bool>(
          valueListenable: logic.phoneVerifyNotifier,
          builder: (_, verifyPhone, __) {
            return ValueListenableBuilder<int>(
              valueListenable: logic.smsCountdownNotifier,
              builder: (_, countdown, __) {
                final text = countdown > 0 ? "验证码已发送($countdown)" : "获取验证码";
                return TextButton(
                  onPressed: verifyPhone && countdown == 0
                      ? () => logic.sendSMS(context)
                      : null,
                  child: LoadingView(
                    loading: countdown == -1,
                    child: Text(text),
                  ),
                );
              },
            );
          },
        ),
      ),
      validator: (v) {
        if (v == null || v.isEmpty) {
          return "验证码不能为空";
        }
        return null;
      },
    );
  }

  @override
  void dispose() {
    logic.dispose();
    super.dispose();
  }
}

/*
* 授权页逻辑处理
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
  final authStateNotifier = ValueChangeNotifier(false);

  // 授权信息保存
  Future<void> authSaved(BuildContext context) async {
    final currState = formKey.currentState;
    if (currState == null || !currState.validate()) return;
    try {
      authStateNotifier.setValue(true);
      final phoneNumber = phoneController.text;
      final code = smsCodeController.text;
      await authApi.auth(phoneNumber: phoneNumber, code: code);
      // 登录成功跳转首页
      routerManage.pushReplacementNamed(RoutePath.home);
    } catch (e) {
      SnackTool.showMessage(context,
          message: debugMode ? e.toString() : "请求授权失败");
      authStateNotifier.setValue(false);
    }
  }

  // 记录短信验证码获取倒计时
  final smsCountdownNotifier = ValueChangeNotifier<int>(0);

  // 发送短信验证码
  Future<void> sendSMS(BuildContext context) async {
    try {
      smsCountdownNotifier.setValue(-1);
      final phoneNumber = phoneController.value.text;
      if (await authApi.sendSMS(phoneNumber: phoneNumber)) {
        _startSmsCountdown();
        if (debugMode) {
          // 开发模式，不会发送验证码，默认使用手机号后四位
          smsCodeController.text =
              phoneNumber.substring(phoneNumber.length - 4);
        }
      }
    } catch (e) {
      SnackTool.showMessage(context, message: "短信验证码发送失败");
      smsCountdownNotifier.setValue(0);
    }
  }

  // 计时器
  Timer? smsCountdownTimer;

  // 短信验证码获取倒计时
  void _startSmsCountdown() {
    final countDown = debugMode ? 5 : 60;
    smsCountdownNotifier.setValue(countDown);
    smsCountdownTimer = Timer.periodic(
      const Duration(seconds: 1),
      (t) {
        final v = smsCountdownNotifier.value;
        smsCountdownNotifier.setValue(v - 1);
        if (smsCountdownNotifier.value <= 0) t.cancel();
      },
    );
  }

  @override
  void dispose() {
    // 销毁控制器和计时器
    smsCountdownTimer?.cancel();
    phoneController.dispose();
    smsCodeController.dispose();
    phoneVerifyNotifier.dispose();
    smsCountdownNotifier.dispose();
    super.dispose();
  }
}
