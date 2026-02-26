#include <jni.h>

extern "C" JNIEXPORT jboolean JNICALL
Java_com_xshare_corebridge_CoreBridge_nativeStartForwarding(
    JNIEnv* /* env */,
    jobject /* thiz */
) {
    return JNI_TRUE;
}

extern "C" JNIEXPORT void JNICALL
Java_com_xshare_corebridge_CoreBridge_nativeStopForwarding(
    JNIEnv* /* env */,
    jobject /* thiz */
) {
}
