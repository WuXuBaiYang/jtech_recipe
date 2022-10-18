import 'package:client/common/model.dart';
import 'package:client/model/model.dart';

/*
* 标签对象
* @author wuxubaiyang
* @Time 2022/10/14 14:28
*/
class TagModel extends BaseModel with BasePart, CreatorPart {
  // 父值
  final String pCode;

  // 值
  final String code;

  // 标签内容
  final String tag;

  // 标签描述
  final String info;

  TagModel.from(obj)
      : pCode = obj?["pCode"] ?? "",
        code = obj?["code"] ?? "",
        tag = obj?["tag"] ?? obj?["name"] ?? "",
        info = obj?["info"] ?? "" {
    initBasePart(obj);
    initCreatorPart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        "pCode": pCode,
        "code": code,
        "tag": tag,
        "info": info,
      };

  // 导出添加标签信息结构
  Map<String, dynamic> toAddInfo() => {
        "tag": tag,
        "info": info,
      };
}
