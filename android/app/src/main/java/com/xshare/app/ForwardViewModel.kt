package com.xshare.app

import androidx.lifecycle.ViewModel
import com.xshare.corebridge.Bridge

class ForwardViewModel(
    private val bridge: Bridge
) : ViewModel() {
    enum class State {
        Idle,
        Running,
        Error
    }

    data class StateHolder<T>(var value: T)

    val state = StateHolder(State.Idle)

    fun startForward() {
        state.value = if (bridge.startForwarding()) State.Running else State.Error
    }

    fun stopForward() {
        bridge.stopForwarding()
        state.value = State.Idle
    }
}
