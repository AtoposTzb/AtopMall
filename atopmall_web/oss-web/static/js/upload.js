// MinIO 预签名 PUT 直传（替代阿里云 OSS 的 POST + policy/signature 方式）

// 缓存签名信息
var g_upload_url = "";
var g_file_url = "";
var g_expire = 0;

// 获取上传凭证（调用后端 /oss/token 接口）
function send_request() {
  var xmlhttp = null;
  if (window.XMLHttpRequest) {
    xmlhttp = new XMLHttpRequest();
  } else if (window.ActiveXObject) {
    xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
  }

  if (xmlhttp != null) {
    // 修改为你自己的 oss-web 服务地址
    serverUrl = "http://192.168.1.6:8083/oss/v1/oss/token";

    xmlhttp.open("GET", serverUrl, false);
    xmlhttp.send(null);
    return xmlhttp.responseText;
  } else {
    alert("Your browser does not support XMLHTTP.");
  }
}

// 检查签名是否过期，过期则重新获取
function get_signature() {
  var now = Date.parse(new Date()) / 1000;
  if (g_expire < now + 3) {
    var body = send_request();
    var obj = JSON.parse(body);
    if (obj.code != 200) {
      alert("获取上传凭证失败: " + obj.msg);
      return false;
    }
    g_upload_url = obj.data.upload_url;
    g_file_url = obj.data.url;
    g_expire = now + 3600; // 签名有效期 1 小时
    return true;
  }
  return true;
}

// 随机文件名
function random_string(len) {
  len = len || 32;
  var chars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678";
  var maxPos = chars.length;
  var pwd = "";
  for (var i = 0; i < len; i++) {
    pwd += chars.charAt(Math.floor(Math.random() * maxPos));
  }
  return pwd;
}

function get_suffix(filename) {
  var pos = filename.lastIndexOf(".");
  var suffix = "";
  if (pos != -1) {
    suffix = filename.substring(pos);
  }
  return suffix;
}

// 计算上传对象名
function calculate_object_name(filename, type) {
  if (type == "local_name") {
    return filename;
  } else if (type == "random_name") {
    return random_string(10) + get_suffix(filename);
  }
  return filename;
}

// 设置上传参数并执行上传
function set_upload_param(up, filename, ret) {
  if (ret == false) {
    ret = get_signature();
  }
  if (!ret) return;

  // 根据单选按钮决定文件名策略
  var type = "local_name";
  var radios = document.getElementsByName("myradio");
  for (var i = 0; i < radios.length; i++) {
    if (radios[i].checked) {
      type = radios[i].value;
      break;
    }
  }

  var objectName = calculate_object_name(filename, type);
  // 替换预签名 URL 中的文件名（PresignedPutObject 的 URL 已包含路径，直接 PUT 即可）
  // MinIO PresignedPutObject 返回的 URL 已经包含了完整路径，直接上传即可

  up.setOption({
    url: g_upload_url,
    multipart: false, // MinIO PUT 直传，不使用 multipart
    headers: {
      "Content-Type": "", // 让浏览器自动检测
    },
  });

  up.start();
}

var uploader = new plupload.Uploader({
  runtimes: "html5,flash,silverlight,html4",
  browse_button: "selectfiles",
  container: document.getElementById("container"),
  flash_swf_url: "lib/plupload-2.1.2/js/Moxie.swf",
  silverlight_xap_url: "lib/plupload-2.1.2/js/Moxie.xap",
  url: "", // 实际 URL 在 BeforeUpload 中动态设置

  filters: {
    mime_types: [
      { title: "Image files", extensions: "jpg,gif,png,bmp" },
      { title: "Zip files", extensions: "zip,rar" },
    ],
    max_file_size: "10mb",
    prevent_duplicates: true,
  },

  init: {
    PostInit: function () {
      document.getElementById("ossfile").innerHTML = "";
      document.getElementById("postfiles").onclick = function () {
        set_upload_param(uploader, "", false);
        return false;
      };
    },

    FilesAdded: function (up, files) {
      plupload.each(files, function (file) {
        document.getElementById("ossfile").innerHTML +=
          '<div id="' +
          file.id +
          '">' +
          file.name +
          " (" +
          plupload.formatSize(file.size) +
          ")<b></b>" +
          '<div class="progress"><div class="progress-bar" style="width: 0%"></div></div>' +
          "</div>";
      });
    },

    BeforeUpload: function (up, file) {
      set_upload_param(up, file.name, false);
    },

    UploadProgress: function (up, file) {
      var d = document.getElementById(file.id);
      d.getElementsByTagName("b")[0].innerHTML =
        "<span>" + file.percent + "%</span>";
      var prog = d.getElementsByTagName("div")[0];
      var progBar = prog.getElementsByTagName("div")[0];
      progBar.style.width = 2 * file.percent + "px";
      progBar.setAttribute("aria-valuenow", file.percent);
    },

    FileUploaded: function (up, file, info) {
      if (info.status == 200) {
        document
          .getElementById(file.id)
          .getElementsByTagName("b")[0].innerHTML =
          "上传成功！访问地址: " + g_file_url;
      } else {
        document
          .getElementById(file.id)
          .getElementsByTagName("b")[0].innerHTML =
          "上传失败，状态码: " + info.status + "，响应: " + info.response;
      }
    },

    Error: function (up, err) {
      if (err.code == -600) {
        document
          .getElementById("console")
          .appendChild(document.createTextNode("\n选择的文件太大了"));
      } else if (err.code == -601) {
        document
          .getElementById("console")
          .appendChild(document.createTextNode("\n选择的文件后缀不对"));
      } else if (err.code == -602) {
        document
          .getElementById("console")
          .appendChild(document.createTextNode("\n这个文件已经上传过一遍了"));
      } else {
        document
          .getElementById("console")
          .appendChild(document.createTextNode("\nError: " + err.response));
      }
    },
  },
});

uploader.init();
