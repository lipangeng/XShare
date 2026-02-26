package com.xshare.corebridge

interface Bridge {
    fun startForwarding(): Boolean
    fun stopForwarding()
}

object CoreBridge : Bridge {
    private val loadFailure: Throwable?
    private val nativeLoaded: Boolean

    init {
        val failure = try {
            System.loadLibrary("corebridge")
            null
        } catch (t: Throwable) {
            t
        }
        loadFailure = failure
        nativeLoaded = failure == null
    }

    private fun ensureLoaded() {
        if (!nativeLoaded) {
            throw IllegalStateException("corebridge native library is not loaded", loadFailure)
        }
    }

    override fun startForwarding(): Boolean {
        ensureLoaded()
        return nativeStartForwarding()
    }

    override fun stopForwarding() {
        ensureLoaded()
        nativeStopForwarding()
    }

    private external fun nativeStartForwarding(): Boolean
    private external fun nativeStopForwarding()
}
