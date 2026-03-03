#include <jni.h>

extern "C" {

JNIEXPORT jobject JNICALL
Java_com_xshare_app_CoreBridgeImpl_startForwarding(JNIEnv *env, jobject thiz) {
    jclass resultClass = env->FindClass("com/xshare/app/Result$Success");
    if (resultClass == nullptr) {
        jclass errorClass = env->FindClass("com/xshare/app/Result$Error");
        jmethodID errorCtor = env->GetMethodID(errorClass, "<init>", "(Ljava/lang/String;)V");
        jstring msg = env->NewStringUTF("Not implemented yet");
        return env->NewObject(errorClass, errorCtor, msg);
    }
    jmethodID successCtor = env->GetMethodID(resultClass, "<init>", "(Ljava/lang/Object;)V");
    jclass unitClass = env->FindClass("kotlin/Unit");
    jfieldID unitField = env->GetStaticFieldID(unitClass, "INSTANCE", "Lkotlin/Unit;");
    jobject unit = env->GetStaticObjectField(unitClass, unitField);
    return env->NewObject(resultClass, successCtor, unit);
}

JNIEXPORT jobject JNICALL
Java_com_xshare_app_CoreBridgeImpl_stopForwarding(JNIEnv *env, jobject thiz) {
    jclass resultClass = env->FindClass("com/xshare/app/Result$Success");
    jmethodID successCtor = env->GetMethodID(resultClass, "<init>", "(Ljava/lang/Object;)V");
    jclass unitClass = env->FindClass("kotlin/Unit");
    jfieldID unitField = env->GetStaticFieldID(unitClass, "INSTANCE", "Lkotlin/Unit;");
    jobject unit = env->GetStaticObjectField(unitClass, unitField);
    return env->NewObject(resultClass, successCtor, unit);
}

JNIEXPORT jobject JNICALL
Java_com_xshare_app_CoreBridgeImpl_getStats(JNIEnv *env, jobject thiz) {
    jclass statsClass = env->FindClass("com/xshare/app/ForwardStats");
    jmethodID statsCtor = env->GetMethodID(statsClass, "<init>", "(JJJJII)V");
    return env->NewObject(statsClass, statsCtor, (jlong)0, (jlong)0, (jlong)0, (jlong)0, (jint)0, (jint)0);
}

JNIEXPORT jboolean JNICALL
Java_com_xshare_app_CoreBridgeImpl_isConnected(JNIEnv *env, jobject thiz) {
    return JNI_TRUE;
}

}