package com.xshare.app

import com.xshare.corebridge.Bridge
import org.junit.Assert.assertEquals
import org.junit.Test

class ForwardViewModelTest {
    @Test
    fun startForward_updatesStateToRunning() {
        val vm = ForwardViewModel(FakeBridge(startResult = true))

        vm.startForward()

        assertEquals(ForwardViewModel.State.Running, vm.state.value)
    }

    @Test
    fun startForward_whenBridgeFails_setsErrorState() {
        val vm = ForwardViewModel(FakeBridge(startResult = false))

        vm.startForward()

        assertEquals(ForwardViewModel.State.Error, vm.state.value)
    }

    @Test
    fun startForward_whenBridgeThrows_setsErrorState() {
        val vm = ForwardViewModel(ThrowingBridge())

        vm.startForward()

        assertEquals(ForwardViewModel.State.Error, vm.state.value)
    }

    @Test
    fun stopForward_callsBridgeStop_andReturnsToStopped() {
        val bridge = FakeBridge(startResult = true)
        val vm = ForwardViewModel(bridge)
        vm.startForward()

        vm.stopForward()

        assertEquals(1, bridge.stopCalls)
        assertEquals(ForwardViewModel.State.Stopped, vm.state.value)
    }

    private class FakeBridge(
        private val startResult: Boolean
    ) : Bridge {
        var stopCalls: Int = 0

        override fun startForwarding(): Boolean = startResult

        override fun stopForwarding() {
            stopCalls += 1
        }
    }

    private class ThrowingBridge : Bridge {
        override fun startForwarding(): Boolean {
            throw IllegalStateException("native unavailable")
        }

        override fun stopForwarding() = Unit
    }
}
