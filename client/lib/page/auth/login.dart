import 'package:client/tool/tool.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/*
* 登录页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<StatefulWidget> createState() => _LoginPageState();
}

/*
* 登录页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _LoginPageState extends State<LoginPage> {
  // 表单key
  final formKey = GlobalKey<FormState>();

  // 手机号输入框控制器
  final phoneController = TextEditingController();

  // 验证码输入框控制器
  final smsCodeController = TextEditingController();

  // 密码输入框控制器
  final passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        padding: const EdgeInsets.all(35),
        alignment: Alignment.center,
        child: Form(
          key: formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              TextFormField(
                controller: phoneController,
                keyboardType: TextInputType.phone,
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                ],
                decoration: const InputDecoration(
                  prefixIcon: Icon(Icons.phone),
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
              ),
              TextFormField(
                controller: smsCodeController,
                keyboardType: TextInputType.number,
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                ],
                decoration: const InputDecoration(prefixIcon: Icon(Icons.sms)),
                validator: (v) {
                  if (v == null || v.isEmpty) {
                    return "短信验证码不能为空";
                  }
                  return null;
                },
              ),
              TextFormField(
                controller: passwordController,
                keyboardType: TextInputType.visiblePassword,
                inputFormatters: [
                  FilteringTextInputFormatter.allow(RegExp(r'[a-zA-Z0-9]')),
                ],
                decoration: const InputDecoration(
                  prefixIcon: Icon(Icons.password),
                ),
                validator: (v) {
                  if (v == null || v.isEmpty) {
                    return "密码不能为空";
                  }
                  return null;
                },
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        child: const Icon(Icons.done),
        onPressed: () {
          formKey.currentState?.validate();
        },
      ),
    );
  }
}
