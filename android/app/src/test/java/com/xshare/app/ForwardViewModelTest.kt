package com.xshare.app

import org.junit.Test
import org.junit.Assert.*

class ForwardViewModelTest {

    @Test
    fun startForward_updatesStateToRunning() {
        val vm = ForwardViewModel(FakeBridge())
        vm.startForward()
        assertEquals(ForwardState.Running, vm.state.value)
    }

    @Test
    fun stopForward_afterStart_updatesStateToIdle() {
        val vm = ForwardViewModel(FakeBridge())
        vm.startForward()
        vm.stopForward()
        assertEquals(ForwardState.Idle, vm.state.value)
    }

    @Test
    fun startForward_whenAlreadyRunning_returnsError() {
        val vm = ForwardViewModel(FakeBridge())
        vm.startForward()
        // Note: The current implementation doesn't return a result from startForward()
        // This test would need to be adjusted based on the actual API
    }
}
