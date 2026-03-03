#include <jni.h>
#include <string>

extern "C" {

JNIEXPORT jboolean JNICALL
Java_com_xshare_app_CoreBridge_nativeIsConnected(JNIEnv *env, jobject thiz) {
    return JNI_FALSE;
}

JNIEXPORT jstring JNICALL
Java_com_xshare_app_CoreBridge_nativeGetVersion(JNIEnv *env, jobject thiz) {
    return env->NewStringUTF("1.0.0");
}

} // extern "C"
