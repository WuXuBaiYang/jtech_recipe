import 'package:client/tool/tool.dart';
import 'package:flutter/material.dart';

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
  // 分页控制器
  final pageController = PageController();

  // 手机号输入框控制器
  final phoneController = TextEditingController();

  // 密码输入框控制器
  final passwordController = TextEditingController();

  // 验证码输入框控制器
  final smsCodeController = TextEditingController();

  // 分页组件保存
  late List<Widget> pageChildren = [
    _buildLoginPage(),
    _buildRegisterPage(),
  ];

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: PageView(
        physics: const NeverScrollableScrollPhysics(),
        controller: pageController,
        children: [
          _buildLoginPage(),
          _buildRegisterPage(),
        ],
        // pageChildren,
      ),
    );
  }

  // 构建登录分页
  Widget _buildLoginPage() {
    return Scaffold(
      body: Container(
        padding: const EdgeInsets.all(35),
        alignment: Alignment.center,
        child: Form(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              TextFormField(
                controller: phoneController,
                validator: (v) {
                  if (v == null || v.isEmpty) {
                    return "手机号不能为空";
                  }
                  if (!Tool.verifyPhone(v)) {
                    return "手机号校验失败";
                  }
                  return null;
                },
                decoration: const InputDecoration(
                  label: Text("手机号"),
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // 构建注册分页
  Widget _buildRegisterPage() {
    return Scaffold(
      body: Text("aaa"),
    );
  }
}
