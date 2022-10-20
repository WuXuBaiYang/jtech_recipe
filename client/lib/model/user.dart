import 'package:client/common/model.dart';
import 'package:client/manage/tag.dart';
import 'package:client/model/tag.dart';
import 'package:client/tool/date.dart';
import 'package:flutter/cupertino.dart';

import 'model.dart';

/*
* 用户信息
* @author wuxubaiyang
* @Time 2022/9/12 19:10
*/
class UserModel extends BaseModel with BasePart {
  // 手机号
  final String phoneNumber;

  // 昵称
  final String nickName;

  // 头像
  final String avatar;

  // 简介
  final String bio;

  // 职业
  final String profession;

  // 性别
  final String genderCode;

  // 生日
  final DateTime birth;

  // 勋章列表
  final List<MedalModel> medals;

  // 个人评价
  final String evaluateCode;

  // 偏好菜系
  final List<String> recipeCuisineCodes;

  // 偏好口味
  final List<String> recipeTasteCodes;

  // 用户经验
  final num exp;

  // 用户等级
  final num level;

  // 当前等级已获得经验
  final num levelExp;

  // 升级所需经验
  final num updateExp;

  // 获取性别
  Future<String> getGender(BuildContext context) async {
    final result = await tagManage.findTag(
      context,
      source: TagSource.userGender,
      code: genderCode,
    );
    return result?.tag ?? '未知';
  }

  // 获取个人评价
  Future<String> getEvaluate(BuildContext context) async {
    final result = await tagManage.findTag(
      context,
      source: TagSource.userEvaluate,
      code: evaluateCode,
    );
    return result?.tag ?? '';
  }

  UserModel.from(obj)
      : phoneNumber = obj?['phoneNumber'] ?? '',
        nickName = obj?['nickName'] ?? '',
        avatar = obj?['avatar'] ?? '',
        bio = obj?['bio'] ?? '',
        profession = obj?['profession'] ?? '',
        genderCode = obj?['genderCode'] ?? '',
        birth = DateTool.parseDate(obj['birth'] ?? '') ?? DateTime(0),
        medals = (obj?['medals'] ?? [])
            .map<MedalModel>((e) => MedalModel.from(e))
            .toList(),
        evaluateCode = obj?['evaluateCode'] ?? '',
        recipeCuisineCodes = (obj?['recipeCuisineCodes'] ?? [])
            .map<String>((e) => '$e')
            .toList(),
        recipeTasteCodes =
            (obj?['recipeTasteCodes'] ?? []).map<String>((e) => '$e').toList(),
        exp = obj?['exp'] ?? 0,
        level = obj?['level'] ?? 0,
        levelExp = obj?['levelExp'] ?? 0,
        updateExp = obj?['updateExp'] ?? 0 {
    initBasePart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        'phoneNumber': phoneNumber,
        'nickName': nickName,
        'avatar': avatar,
        'bio': bio,
        'profession': profession,
        'genderCode': genderCode,
        'birth': birth.toString(),
        'medals': medals.map((e) => e.to()).toList(),
        'evaluateCode': evaluateCode,
        'recipeCuisineCodes': recipeCuisineCodes,
        'recipeTasteCodes': recipeTasteCodes,
        'exp': exp,
        'level': level,
        'levelExp': levelExp,
        'updateExp': updateExp,
      };

  // 获取编辑结构
  Map<String, dynamic> toModifyInfo() => {
        'nickName': nickName,
        'avatar': avatar,
        'bio': bio,
        'profession': profession,
        'genderCode': genderCode,
        'birth': birth.toString(),
        'evaluateCode': evaluateCode,
        'recipeCuisineCodes': recipeCuisineCodes,
        'recipeTasteCodes': recipeTasteCodes,
      };
}

/*
* 勋章对象
* @author wuxubaiyang
* @Time 2022/10/14 14:16
*/
class MedalModel extends BaseModel with BasePart {
  // 图标
  final String logo;

  // 名称
  final String name;

  // 稀有度
  final String rarityCode;

  MedalModel.from(obj)
      : logo = obj?['logo'] ?? '',
        name = obj?['name'] ?? '',
        rarityCode = obj?['rarityCode'] ?? '' {
    initBasePart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        'logo': logo,
        'name': name,
        'rarityCode': rarityCode,
      };

  // 获取编辑结构
  Map<String, dynamic> toModifyInfo() => {
        'logo': logo,
        'name': name,
        'rarityCode': rarityCode,
      };
}

/*
* 用户收货地址
* @author wuxubaiyang
* @Time 2022/10/14 14:26
*/
class UserAddressModel extends BaseModel with BasePart, CreatorPart {
  // 收货人
  final String receiver;

  // 联系方式
  final String contact;

  // 地址字典码
  final List<String> addressCodes;

  // 地址详情
  final String addressDetail;

  // 标签字典码
  final String tagCode;

  // 标签
  final TagModel? tag;

  // 是否为默认地址
  final bool isDefault;

  // 排序
  final num order;

  UserAddressModel.from(obj)
      : receiver = obj?['receiver'] ?? '',
        contact = obj?['contact'] ?? '',
        addressCodes =
            (obj?['addressCodes'] ?? []).map<String>((e) => '$e').toList(),
        addressDetail = obj?['addressDetail'] ?? '',
        tagCode = obj?['tagCode'] ?? '',
        tag = obj?['tag'] != null ? TagModel.from(obj?['tag'] ?? {}) : null,
        isDefault = obj?['default'] ?? false,
        order = obj?['order'] ?? 0 {
    initBasePart(obj);
    initCreatorPart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        'receiver': receiver,
        'contact': contact,
        'addressCodes': addressCodes,
        'addressDetail': addressDetail,
        'tag': tag?.to(),
        'tagCode': tagCode,
        'default': isDefault,
        'order': order,
      };

  // 获取编辑结构
  Map<String, dynamic> toModifyInfo() => {
        'receiver': receiver,
        'contact': contact,
        'addressCodes': addressCodes,
        'addressDetail': addressDetail,
        'tagCode': tagCode,
        'default': isDefault,
        'order': order,
      };
}
