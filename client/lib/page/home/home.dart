import 'dart:convert';

import 'package:client/model/tag.dart';
import 'package:flutter/material.dart';
import 'package:json_path/json_path.dart';

/*
* 首页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePageState();
}

/*
* 首页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("首页"),
      ),
      body: Center(
        child: TextButton(
          child: Text("点击测试"),
          onPressed: () async {

          },
        ),
      ),
      floatingActionButton: FloatingActionButton(
        child: const Icon(Icons.device_hub),
        onPressed: () async {},
      ),
    );
  }
}
