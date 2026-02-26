package com.xshare.app

import androidx.lifecycle.ViewModel
import com.xshare.corebridge.Bridge

class ForwardViewModel(
    private val bridge: Bridge
) : ViewModel() {
    enum class State {
        Stopped,
        Running,
        Error
    }

    data class StateHolder<T>(var value: T)

    val state = StateHolder(State.Stopped)

    fun startForward() {
        state.value = try {
            if (bridge.startForwarding()) State.Running else State.Error
        } catch (_: Throwable) {
            State.Error
        }
    }

    fun stopForward() {
        state.value = try {
            bridge.stopForwarding()
            State.Stopped
        } catch (_: Throwable) {
            State.Error
        }
    }
}
