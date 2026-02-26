package com.xshare.corebridge

interface Bridge {
    fun startForwarding(): Boolean
    fun stopForwarding()
}

object CoreBridge : Bridge {
    init {
        runCatching {
            System.loadLibrary("corebridge")
        }
    }

    override fun startForwarding(): Boolean = nativeStartForwarding()

    override fun stopForwarding() {
        nativeStopForwarding()
    }

    private external fun nativeStartForwarding(): Boolean
    private external fun nativeStopForwarding()
}
