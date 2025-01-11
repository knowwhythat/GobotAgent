import os
import platform
import re


def is_macos():
    os_type = platform.system()
    if os_type.startswith("Darwin"):
        return True
    return False


def is_linux():
    os_type = platform.system()
    if os_type.startswith("Linux"):
        return True
    return False


def is_win32():
    os_type = platform.system()
    if os_type.startswith("Windows"):
        return True
    return False


def get_browser_path():
    chrome_path = get_chrome_location()
    if chrome_path:
        return chrome_path
    edge_path = get_edge_path()
    if edge_path:
        return edge_path
    raise Exception("未找到浏览器")


def get_chrome_location():
    if is_win32():
        import winreg

        # 尝试从注册表中获取Chrome的安装路径
        try:
            # 64位系统
            reg_path = r"SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe"
            with winreg.OpenKey(
                    winreg.HKEY_LOCAL_MACHINE, reg_path, 0, winreg.KEY_READ
            ) as reg_key:
                chrome_path, _ = winreg.QueryValueEx(reg_key, "")
                return chrome_path
        except FileNotFoundError:
            pass
        try:
            # 32位系统
            reg_path = r"SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe"
            with winreg.OpenKey(
                    winreg.HKEY_LOCAL_MACHINE, reg_path, 0, winreg.KEY_READ
            ) as reg_key:
                chrome_path, _ = winreg.QueryValueEx(reg_key, "")
                return chrome_path
        except FileNotFoundError:
            pass

            # 如果注册表中没有找到，尝试在默认安装路径中查找
        default_paths = [
            r"C:\Program Files\Google\Chrome\Application\chrome.exe",
            r"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe",
            r"D:\Program Files\Google\Chrome\Application\chrome.exe",
            r"D:\Program Files (x86)\Google\Chrome\Application\chrome.exe",
            os.path.join(
                os.path.expanduser("~"),
                "AppData",
                "Local",
                "Google",
                "Chrome",
                "Application",
                "chrome.exe",
            ),
        ]
        for path in default_paths:
            if os.path.exists(path):
                return path

        return None
    elif is_linux():
        default_paths = [
            r"/opt/apps/cn.google.chrome-pre/files/google/chrome/chrome",
            r"/opt/apps/cn.google.chrome/files/google/chrome/chrome",
            r"/opt/google/chrome/chrome",
            r"/opt/google/chrome-pre/chrome",
        ]
        for path in default_paths:
            if os.path.exists(path):
                return path
        default_path = [
            "/var/lib/linglong/layers/main/cn.google.chrome.linyaps",
            "/var/lib/linglong/layers/main/cn.google.chrome-pre.linyaps",
        ]
        version_pattern = re.compile(r"\d+.\d+.\d+.\d+")
        for path in default_path:
            if os.path.exists(path):
                for temp in os.listdir(path):
                    if version_pattern.match(temp):
                        result = os.path.join(
                            path,
                            temp,
                            platform.machine(),
                            r"binary/files/google/chrome/chrome",
                        )
                        if os.path.exists(result):
                            return result
        return None


def get_edge_path_from_registry():
    import winreg

    try:
        # 尝试从注册表中获取Edge的安装路径
        with winreg.OpenKey(
                winreg.HKEY_LOCAL_MACHINE,
                r"SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe",
        ) as key:
            edge_path, _ = winreg.QueryValueEx(key, "")
            return edge_path
    except FileNotFoundError:
        # 如果在64位系统上，还需要检查32位注册表路径
        try:
            with winreg.OpenKey(
                    winreg.HKEY_LOCAL_MACHINE,
                    r"SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe",
            ) as key:
                edge_path, _ = winreg.QueryValueEx(key, "")
                return edge_path
        except FileNotFoundError:
            return None


def get_edge_path_from_env():
    # 尝试从环境变量中获取Edge的安装路径
    for path in os.environ["PATH"].split(os.pathsep):
        edge_path = os.path.join(path, "msedge.exe")
        if os.path.exists(edge_path):
            return os.path.dirname(edge_path)
    return None


def get_edge_path():
    # 首先尝试从注册表获取路径
    edge_path = get_edge_path_from_registry()
    if edge_path:
        return edge_path

        # 如果注册表中没有找到，尝试从环境变量中获取路径
    edge_path = get_edge_path_from_env()
    if edge_path:
        return edge_path

        # 如果都没有找到，尝试一些常见的默认安装路径
    default_paths = [
        r"C:\Program Files\Microsoft\Edge\Application\msedge.exe",
        r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe",
    ]
    for path in default_paths:
        if os.path.exists(path):
            return path

    return None


def get_default_browser_path():
    default_paths = [
        r"/usr/bin/browser",
    ]
    for path in default_paths:
        if os.path.exists(path):
            return path
    deepin_browser_path = r"/var/lib/linglong/layers/main/org.deepin.browser"
    version_pattern = re.compile(r"\d+.\d+.\d+.\d+")
    if os.path.exists(deepin_browser_path):
        for temp in os.listdir(deepin_browser_path):
            if version_pattern.match(temp):
                result = os.path.join(
                    deepin_browser_path,
                    temp,
                    platform.machine(),
                    r"binary/files/share/browser/browser",
                )
                if os.path.exists(result):
                    return result

    return None