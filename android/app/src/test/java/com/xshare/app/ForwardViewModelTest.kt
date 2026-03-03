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
        val result = vm.startForward()
        assertTrue(result.isError)
    }
}
