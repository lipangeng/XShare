package com.xshare.app

import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

sealed class ForwardState {
    object Idle : ForwardState()
    object Starting : ForwardState()
    object Running : ForwardState()
    object Stopping : ForwardState()
    data class Error(val message: String) : ForwardState()
}

data class ForwardStats(
    val uplinkPackets: Long = 0,
    val uplinkBytes: Long = 0,
    val downlinkPackets: Long = 0,
    val downlinkBytes: Long = 0,
    val activeTcpSessions: Int = 0,
    val activeUdpSessions: Int = 0
)

sealed class Result<out T> {
    data class Success<T>(val data: T) : Result<T>()
    data class Error(val message: String) : Result<Nothing>()
    
    val isSuccess: Boolean get() = this is Success
    val isError: Boolean get() = this is Error
    
    fun getOrNull(): T? = (this as? Success)?.data
    fun errorOrNull(): String? = (this as? Error)?.message
}

interface CoreBridge {
    fun startForwarding(): Result<Unit>
    fun stopForwarding(): Result<Unit>
    fun getStats(): Result<ForwardStats>
    fun isConnected(): Boolean
}

class ForwardViewModel(private val bridge: CoreBridge) {
    
    private val _state = MutableStateFlow<ForwardState>(ForwardState.Idle)
    val state: StateFlow<ForwardState> = _state.asStateFlow()
    
    private val _stats = MutableStateFlow(ForwardStats())
    val stats: StateFlow<ForwardStats> = _stats.asStateFlow()
    
    fun startForward() {
        if (_state.value != ForwardState.Idle) {
            return
        }
        
        _state.value = ForwardState.Starting
        
        when (val result = bridge.startForwarding()) {
            is Result.Success -> {
                _state.value = ForwardState.Running
            }
            is Result.Error -> {
                _state.value = ForwardState.Error(result.message)
            }
        }
    }
    
    fun stopForward() {
        if (_state.value != ForwardState.Running) {
            return
        }
        
        _state.value = ForwardState.Stopping
        
        when (val result = bridge.stopForwarding()) {
            is Result.Success -> {
                _state.value = ForwardState.Idle
            }
            is Result.Error -> {
                _state.value = ForwardState.Error(result.message)
            }
        }
    }
    
    fun refreshStats() {
        when (val result = bridge.getStats()) {
            is Result.Success -> {
                _stats.value = result.data
            }
            is Result.Error -> {
                // Keep current stats on error
            }
        }
    }
}

class FakeBridge : CoreBridge {
    private var running = false
    
    override fun startForwarding(): Result<Unit> {
        if (running) {
            return Result.Error("Already running")
        }
        running = true
        return Result.Success(Unit)
    }
    
    override fun stopForwarding(): Result<Unit> {
        if (!running) {
            return Result.Error("Not running")
        }
        running = false
        return Result.Success(Unit)
    }
    
    override fun getStats(): Result<ForwardStats> {
        return Result.Success(ForwardStats())
    }
    
    override fun isConnected(): Boolean = true
}
