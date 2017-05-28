// ./docker/daemon/checkpoint.go   //peter

package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/docker/docker/api/types"
)

// ContainerCheckpoint checkpoints the process running in a container with CRIU
func (daemon *Daemon) ContainerCheckpoint(name string, opts *types.CriuConfig) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}
	if !container.IsRunning() {
		return fmt.Errorf("Container %s not running", name)
	}

	if opts.ImagesDirectory == "" {
		opts.ImagesDirectory = filepath.Join(container.Root, "criu.image")
		if err := os.MkdirAll(opts.ImagesDirectory, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}

	if opts.WorkDirectory == "" {
		opts.WorkDirectory = filepath.Join(container.Root, "criu.work")
		if err := os.MkdirAll(opts.WorkDirectory, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}
        
      
       /// add by peter 
/*   make it empty if empty
       if opts.PrevImagesDirectory == "" {
		opts.PrevImagesDirectory = filepath.Join(container.Root, "criu.preimage")
		if err := os.MkdirAll(opts.PrevImagesDirectory, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}
*/
          
       // add by peter
/*
        if  opts.PageServer == true {
              if opts.Address =="" {
                   return fmt.Errorf("Page Server enabled but address is not assigned");
              }
              if opts.Port =="" || opts.Port == -1 {
                  return fmt.Errorf("Page Server enabled but port is not assigned");
              } 
           /*
           //   p,err := strconv.Atoi(opts.Port)
           //   if err !=nil && p==0{
           //        return fmt.Errorf("Port=%s is not number", opts.Port)
           //   } 
           
        }
      */

	if err := daemon.Checkpoint(container, opts); err != nil {
		return fmt.Errorf("Cannot checkpoint container %s: %s", name, err)
	}

	container.SetCheckpointed(opts.LeaveRunning)
	daemon.LogContainerEvent(container, "checkpoint")

	if opts.LeaveRunning == false {
		daemon.Cleanup(container)
	}

	if err := container.ToDisk(); err != nil {
		return fmt.Errorf("Cannot update config for container: %s", err)
	}

	return nil
}
